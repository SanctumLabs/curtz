// Command curtz runs the Curtz HTTP API.
//
// Only the Identity bounded context is wired up so far; the URL context follows once its
// application layer exists.
package main

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"sync"
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
	"github.com/sanctumlabs/curtz/app/pkg"
	cacheredis "github.com/sanctumlabs/curtz/app/pkg/infra/cache/redis"
	"github.com/sanctumlabs/curtz/app/pkg/infra/database"
	"github.com/sanctumlabs/curtz/app/pkg/infra/database/postgres"
	"github.com/sanctumlabs/curtz/app/pkg/infra/monitoring/health"
	"github.com/sanctumlabs/curtz/app/pkg/infra/server"
	"github.com/sanctumlabs/curtz/app/pkg/infra/server/router"
	"github.com/sanctumlabs/curtz/app/pkg/infra/telemetry"
	"github.com/sanctumlabs/curtz/app/pkg/jwt"
	recoveryutils "github.com/sanctumlabs/curtz/app/pkg/utils/recover"
)

const (
	baseURI = "/api/v1/curtz"

	// redisStartupPingTimeout bounds the one ping that only decides which startup line is logged.
	redisStartupPingTimeout = 2 * time.Second

	// telemetryFlushTimeout bounds the final export of spans and metrics at shutdown.
	telemetryFlushTimeout = 5 * time.Second
)

// healthcheckTimeout bounds the container health probe. It is a variable so a test can shorten it.
var healthcheckTimeout = 2 * time.Second

// setupTelemetry is telemetry.Setup. It is a variable so a test can replace it.
var setupTelemetry = telemetry.Setup

func main() {
	if len(os.Args) > 1 && os.Args[1] == "healthcheck" {
		os.Exit(healthcheck(os.LookupEnv, os.Stdout))
	}

	dotenvErr := godotenv.Load()

	cfg, err := config.Load(os.LookupEnv)
	// The logger is installed once the configuration is read, so LOG_LEVEL and LOG_FORMAT apply. On a configuration error
	// Load still returns the default logging settings, so the error itself is logged as JSON like everything else.
	slog.SetDefault(telemetry.NewLogger(os.Stdout, cfg.Logging.Format, cfg.Logging.Level, telemetry.ServiceName()))
	if dotenvErr != nil {
		slog.Warn("no .env file found, relying on the environment", "error", dotenvErr)
	}
	if err != nil {
		slog.Error("invalid configuration", "error", err)
		os.Exit(1)
	}
	for _, warning := range cfg.Warnings {
		slog.Warn(warning)
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

// healthcheck probes the API running on this host and returns the process exit code: 0 when GET /health answers 200,
// 1 otherwise. It reads only the server settings, so the container's HEALTHCHECK needs no database or Redis settings and
// works in an image that has no shell, curl or wget. It does not read .env: the container's environment is the contract.
func healthcheck(lookup config.Lookup, out io.Writer) int {
	server, err := config.LoadServer(lookup)
	if err != nil {
		_, _ = fmt.Fprintf(out, "healthcheck: invalid configuration: %v\n", err)
		return 1
	}

	target := "http://" + net.JoinHostPort(probeHost(server.Host), strconv.Itoa(server.Port)) + probes.LivePath
	ctx, cancel := context.WithTimeout(context.Background(), healthcheckTimeout)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, target, nil)
	if err != nil {
		_, _ = fmt.Fprintf(out, "healthcheck: %v\n", err)
		return 1
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		_, _ = fmt.Fprintf(out, "healthcheck: %v\n", err)
		return 1
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		_, _ = fmt.Fprintf(out, "healthcheck: %s answered %d\n", target, resp.StatusCode)
		return 1
	}
	_, _ = fmt.Fprintln(out, "healthcheck: ok")
	return 0
}

// probeHost turns a wildcard bind address, which cannot be dialed, into the loopback address that reaches it. Server
// binds exactly the configured host, so any other host is dialed as given.
func probeHost(host string) string {
	switch host {
	case "", "0.0.0.0", "::":
		return "127.0.0.1"
	}
	return host
}

// startTelemetry installs the OpenTelemetry SDK and returns the function that flushes it. Call the flush once the server
// has drained, so the drain's spans and the final metrics are exported; it flushes only once, however often it is called.
// Telemetry never stops the API: when the SDK cannot start the API runs without it, and a failing flush is only logged.
func startTelemetry(ctx context.Context, cfg config.App) (flush func()) {
	shutdown, err := setupTelemetry(ctx, telemetry.Options{ServiceVersion: pkg.Version, Environment: cfg.Environment})
	if err != nil {
		slog.WarnContext(ctx, "telemetry is disabled: the OpenTelemetry SDK could not start", "error", err)
		return func() {}
	}
	var once sync.Once
	return func() {
		once.Do(func() {
			// ctx is already cancelled when this runs (that is what began the shutdown), so the flush gets its own deadline.
			flushCtx, cancel := context.WithTimeout(context.WithoutCancel(ctx), telemetryFlushTimeout)
			defer cancel()
			if err := shutdown(flushCtx); err != nil {
				slog.WarnContext(flushCtx, "flushing telemetry", "error", err)
			}
		})
	}
}

// run builds the API from cfg and serves it until ctx is cancelled, then drains in-flight requests and closes the
// data clients. Postgres is required: failing to reach it is an error. Redis is optional: failing to reach it is
// logged and the API carries on (readiness reports it as down).
func run(ctx context.Context, cfg config.App) error {
	// The deferred flush covers the early returns; the normal path flushes right after the server drains (below).
	flushTelemetry := startTelemetry(ctx, cfg)
	defer flushTelemetry()

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
		ProbePaths:  []string{probes.LivePath, probes.ReadyPath},
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
		},
		PublicPrefixes: []string{"/docs/"},
	}))

	srv.RegisterHandlers([]router.Router{
		probes.NewRouter(registry),
		identityapi.NewRouter(baseURI, identityService),
	})

	serveErr := srv.Serve(ctx, cfg.ShutdownTimeout, registry.SetDraining)

	// The server has drained or given up: export the last spans and metrics now, before the data clients close. Closing the
	// Postgres pool waits for every connection a stuck request still holds, which can outlast the container's stop grace
	// period and would cost the final export of exactly the incident worth diagnosing.
	flushTelemetry()
	return serveErr
}
