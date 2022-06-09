package main

import (
	"fmt"
	"strconv"

	"github.com/joho/godotenv"
	authApi "github.com/sanctumlabs/curtz/app/api/auth"
	"github.com/sanctumlabs/curtz/app/api/health"
	urlApi "github.com/sanctumlabs/curtz/app/api/url"
	"github.com/sanctumlabs/curtz/app/config"
	"github.com/sanctumlabs/curtz/app/internal/core/urlsvc"
	"github.com/sanctumlabs/curtz/app/internal/core/usersvc"
	"github.com/sanctumlabs/curtz/app/internal/repositories"
	"github.com/sanctumlabs/curtz/app/internal/services/auth"
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
	port := env.EnvOr(EnvPort, "8080")
	host := env.EnvOr(EnvDatabaseHost, "localhost")
	database := env.EnvOr(EnvDatabase, "curtzdb")
	databaseUser := env.EnvOr(EnvDatabaseUsername, "curtzUser")
	databasePass := env.EnvOr(EnvDatabasePassword, "curtzPass")
	databasePort := env.EnvOr(EnvDatabasePort, "5432")

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
		Database: config.DatabaseConfig{
			Host:     host,
			Database: database,
			User:     databaseUser,
			Password: databasePass,
			Port:     databasePort,
		},
	}

	srv := server.NewServer(&configuration)

	authService := auth.NewService(configuration.Auth)
	corsMiddleware := middleware.NewCORSMiddleware(configuration.CorsHeaders)
	loggingMiddleware := middleware.NewLoggingMiddleware(configuration.Logging)
	recoveryMiddleware := middleware.NewRecoveryMiddleware()
	authMiddleware := middleware.NewAuthMiddleware(configuration.Auth, authService)

	repository := repositories.NewRepository(configuration.Database)
	urlService := urlsvc.NewUrlSvc(repository.GetUrlRepo())
	userService := usersvc.NewUserSvc(repository.GetUserRepo())

	// setup routers
	routers := []router.Router{
		urlApi.NewUrlRouter(urlService),
		authApi.NewRouter(userService),
		health.NewHealthRouter(),
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
