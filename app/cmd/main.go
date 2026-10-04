// Command curtz runs the Curtz HTTP API.
//
// Only the Identity bounded context is wired up so far; the URL context follows once its
// application layer exists.
package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/joho/godotenv"
	apimiddleware "github.com/sanctumlabs/curtz/app/api/middleware"
	"github.com/sanctumlabs/curtz/app/api/probes"
	identityapi "github.com/sanctumlabs/curtz/app/api/v1/identity"
	"github.com/sanctumlabs/curtz/app/config"
	"github.com/sanctumlabs/curtz/app/internal/adapters/jwtauth"
	"github.com/sanctumlabs/curtz/app/internal/adapters/notifications"
	identitydatastore "github.com/sanctumlabs/curtz/app/internal/adapters/postgres/identity"
	identityapp "github.com/sanctumlabs/curtz/app/internal/application/identity"
	cacheredis "github.com/sanctumlabs/curtz/app/pkg/infra/cache/redis"
	"github.com/sanctumlabs/curtz/app/pkg/infra/database"
	"github.com/sanctumlabs/curtz/app/pkg/infra/database/postgres"
	"github.com/sanctumlabs/curtz/app/pkg/infra/monitoring/health"
	"github.com/sanctumlabs/curtz/app/pkg/infra/server"
	"github.com/sanctumlabs/curtz/app/pkg/infra/server/router"
	"github.com/sanctumlabs/curtz/app/pkg/jwt"
	recoveryutils "github.com/sanctumlabs/curtz/app/pkg/utils/recover"
)

const (
	baseURI = "/api/v1/curtz"

	// redisStartupPingTimeout bounds the one ping that only decides which startup line is logged.
	redisStartupPingTimeout = 2 * time.Second
)

func main() {
	if err := godotenv.Load(); err != nil {
		slog.Warn("no .env file found, relying on the environment", "error", err)
	}

	cfg, err := config.Load(os.LookupEnv)
	if err != nil {
		slog.Error("invalid configuration", "error", err)
		os.Exit(1)
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	go func() {
		<-ctx.Done()
		// Hand signal handling back to the runtime as soon as the first signal arrives, so a second Ctrl-C ends a
		// stuck drain immediately instead of being swallowed.
		stop()
	}()

	err = run(ctx, cfg)
	stop()
	if err != nil {
		slog.Error("curtz stopped", "error", err)
		os.Exit(1)
	}
}

// run builds the API from cfg and serves it until ctx is cancelled, then drains in-flight requests and closes the
// data clients. Postgres is required: failing to reach it is an error. Redis is optional: failing to reach it is
// logged and the API carries on (readiness reports it as down).
func run(ctx context.Context, cfg config.App) error {
	dbClient, err := postgres.NewPostgresClient(cfg.Database.Postgres)
	if err != nil {
		return fmt.Errorf("connect to postgres: %w", err)
	}
	defer dbClient.Close()

	cache, err := cacheredis.NewRedisClient(cfg.Redis)
	if err != nil {
		return fmt.Errorf("create redis client: %w", err)
	}
	defer func() {
		if closeErr := cache.Close(); closeErr != nil {
			slog.Warn("closing redis", "error", closeErr)
		}
	}()

	pingCtx, cancel := context.WithTimeout(ctx, redisStartupPingTimeout)
	if pingErr := cache.Ping(pingCtx); pingErr != nil {
		slog.Warn("redis is down, continuing without it", "addresses", cfg.Redis.Address, "error", pingErr)
	} else {
		slog.Info("redis is up", "addresses", cfg.Redis.Address)
	}
	cancel()

	registry := health.NewRegistry(health.DefaultCheckTimeout)
	registry.Add(health.Check{Name: "postgres", Required: true, Fn: dbClient.HealthCheck})
	registry.Add(health.Check{Name: "redis", Required: false, Fn: cache.Ping})

	dbConfig := database.Config{
		OperationTimeout: cfg.Database.OperationTimeout,
		RetryConfig:      recoveryutils.DefaultRetryConfig,
	}

	tokenService := jwtauth.NewTokenService(cfg.Auth, jwt.New())

	// No email transport is configured yet, so verification links are logged rather than sent.
	notifier := notifications.NewEmailNotifier(cfg.Server.BaseURL, notifications.NewLoggingEmailSender())

	identityService := identityapp.NewService(
		identitydatastore.NewUserDatastoreAdapter(dbClient, dbConfig),
		tokenService,
		notifier,
	)

	srv := server.NewServer(server.ServerConfig{
		Header:      cfg.Server.Header,
		Host:        cfg.Server.Host,
		Port:        cfg.Server.Port,
		AppName:     cfg.Server.Name,
		Version:     cfg.Server.Version,
		Environment: cfg.Environment,
	})

	// Everything is authenticated unless it is listed here. Registration, login, token refresh
	// and email verification must be reachable without a token, by definition, and so must the probes.
	srv.Use(apimiddleware.AuthMiddleware(apimiddleware.AuthConfig{
		TokenService: tokenService,
		PublicPaths: []string{
			baseURI + "/auth/register",
			baseURI + "/auth/login",
			baseURI + "/auth/oauth/token",
			baseURI + "/auth/verify",
			probes.LivePath,
			probes.ReadyPath,
			"/metrics",
		},
		PublicPrefixes: []string{"/docs/"},
	}))

	srv.RegisterHandlers([]router.Router{
		probes.NewRouter(registry),
		identityapi.NewRouter(baseURI, identityService),
	})

	return srv.Serve(ctx, cfg.ShutdownTimeout, registry.SetDraining)
}
