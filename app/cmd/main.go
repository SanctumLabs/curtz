package main

import (
	"fmt"
	"strconv"

	"github.com/joho/godotenv"
	"github.com/sanctumlabs/curtz/app/api/health"
	authApi "github.com/sanctumlabs/curtz/app/api/v1/auth"
	"github.com/sanctumlabs/curtz/app/api/v1/client"
	"github.com/sanctumlabs/curtz/app/api/v1/url"
	"github.com/sanctumlabs/curtz/app/config"
	"github.com/sanctumlabs/curtz/app/internal/core/urlsvc"
	"github.com/sanctumlabs/curtz/app/internal/core/usersvc"
	"github.com/sanctumlabs/curtz/app/internal/repositories"
	"github.com/sanctumlabs/curtz/app/internal/services/auth"
	"github.com/sanctumlabs/curtz/app/internal/services/cache"
	"github.com/sanctumlabs/curtz/app/internal/services/notifications"
	"github.com/sanctumlabs/curtz/app/internal/services/notifications/email"
	"github.com/sanctumlabs/curtz/app/server"
	"github.com/sanctumlabs/curtz/app/server/middleware"
	"github.com/sanctumlabs/curtz/app/server/router"
	"github.com/sanctumlabs/curtz/app/tools/env"
	"github.com/sanctumlabs/curtz/app/tools/logger"
)

const (
	Env                 = "ENV"
	EnvLogLevel         = "LOG_LEVEL"
	EnvLogJsonOutput    = "LOG_JSON_OUTPUT"
	EnvPort             = "PORT"
	EnvDatabaseHost     = "DATABASE_HOST"
	EnvDatabase         = "DATABASE"
	EnvDatabaseUsername = "DATABASE_USERNAME"
	EnvDatabasePassword = "DATABASE_PASSWORD"
	EnvDatabasePort     = "DATABASE_PORT"
	EnvAuthSecret       = "AUTH_SECRET"
	EnvAuthExpireDelta  = "AUTH_EXPIRE_DELTA"
	EnvAuthIssuer       = "AUTH_ISSUER"
	EnvCacheHost        = "CACHE_HOST"
	EnvCacheUsername    = "CACHE_USERNAME"
	EnvCachePassword    = "CACHE_PASSWORD"
	EnvCachePort        = "CACHE_PORT"
	EnvCacheRequireAuth = "CACHE_REQUIRE_AUTH"
)

func main() {
	log := logger.NewLogger("curtz-api")

	err := godotenv.Load()
	if err != nil {
		log.Warn("Error loading .env file. Using defaults")
	}

	environment := env.EnvOr(Env, "development")
	logLevel := env.EnvOr(EnvLogLevel, "debug")
	logJsonOutput := env.EnvOr(EnvLogJsonOutput, "true")
	port := env.EnvOr(EnvPort, "8085")
	host := env.EnvOr(EnvDatabaseHost, "localhost")
	database := env.EnvOr(EnvDatabase, "curtzdb")
	databaseUser := env.EnvOr(EnvDatabaseUsername, "curtzUser")
	databasePass := env.EnvOr(EnvDatabasePassword, "curtzPassword")
	databasePort := env.EnvOr(EnvDatabasePort, "27017")
	authSecret := env.EnvOr(EnvAuthSecret, "curtz-secret")
	authExpireDelta := env.EnvOr(EnvAuthExpireDelta, "6")
	authIssuer := env.EnvOr(EnvAuthIssuer, "curtz")
	cacheHost := env.EnvOr(EnvCacheHost, "localhost")
	cachePort := env.EnvOr(EnvCachePort, "6379")
	cacheUsername := env.EnvOr(EnvCacheUsername, "curtzUser")
	cachePassword := env.EnvOr(EnvCachePassword, "curtzPassword")
	cacheRequireAuth := env.EnvOr(EnvCacheRequireAuth, "false")

	expireDelta, err := strconv.Atoi(authExpireDelta)
	if err != nil {
		expireDelta = 6
	}

	cacheNeedsAuth, err := strconv.ParseBool(cacheRequireAuth)
	if err != nil {
		cacheNeedsAuth = false
	}

	enableJsonOutput, err := strconv.ParseBool(logJsonOutput)
	if err != nil {
		enableJsonOutput = true
	}

	configuration := config.Config{
		Env:  environment,
		Port: port,
		Logging: config.LoggingConfig{
			Level:            logLevel,
			EnableJSONOutput: enableJsonOutput,
		},
		Auth: config.AuthConfig{
			Jwt: config.Jwt{
				Secret:      authSecret,
				ExpireDelta: expireDelta,
				Issuer:      authIssuer,
			},
		},
		Database: config.DatabaseConfig{
			Host:     host,
			Database: database,
			User:     databaseUser,
			Password: databasePass,
			Port:     databasePort,
		},
		Cache: config.CacheConfig{
			Host:        cacheHost,
			Port:        cachePort,
			Username:    cacheUsername,
			Password:    cachePassword,
			RequireAuth: cacheNeedsAuth,
		},
	}

	srv := server.NewServer(&configuration)

	authService := auth.NewService(configuration.Auth)
	corsMiddleware := middleware.NewCORSMiddleware(configuration.CorsHeaders)
	loggingMiddleware := middleware.NewLoggingMiddleware(configuration.Logging)
	recoveryMiddleware := middleware.NewRecoveryMiddleware()
	authMiddleware := middleware.NewAuthMiddleware(configuration.Auth, authService)

	repository := repositories.NewRepository(configuration.Database)
	emailSvc := email.NewEmailSvc()
	notificationSvc := notifications.NewNotificationSvc(emailSvc)
	cache := cache.New(configuration.Cache)

	userService := usersvc.NewUserSvc(repository.GetUserRepo(), notificationSvc)
	urlService := urlsvc.NewUrlSvc(repository.GetUrlRepo(), userService, cache)

	baseUri := "/api/v1/curtz"

	// setup routers
	routers := []router.Router{
		url.NewUrlRouter(baseUri, urlService),
		authApi.NewRouter(baseUri, userService, authService),
		health.NewHealthRouter(),
		client.NewClientRouter(urlService),
	}

	// initialize routers
	srv.InitRouter(routers...)

	// use middlewares
	srv.UseMiddleware(loggingMiddleware)
	srv.UseMiddleware(corsMiddleware)
	srv.UseMiddleware(recoveryMiddleware)
	srv.UseMiddleware(authMiddleware)

	appServer := srv.CreateServer()

	// start & run the server
	err = appServer.Run(fmt.Sprintf(":%s", port))
	if err != nil {
		_, msg := fmt.Printf("Failed to start Server %s", err)
		log.Error(msg)
		panic(msg)
	}
}
