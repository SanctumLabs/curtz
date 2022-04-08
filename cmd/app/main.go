package main

import (
	"fmt"
	"github.com/joho/godotenv"
	"github.com/sanctumlabs/curtz/api/health"
	urlApi "github.com/sanctumlabs/curtz/api/url"
	"github.com/sanctumlabs/curtz/config"
	"github.com/sanctumlabs/curtz/internal/core/usecases/url"
	"github.com/sanctumlabs/curtz/internal/repositories"
	"github.com/sanctumlabs/curtz/internal/services/urlsvc"
	"github.com/sanctumlabs/curtz/server"
	"github.com/sanctumlabs/curtz/server/middleware"
	"github.com/sanctumlabs/curtz/server/router"
	"github.com/sanctumlabs/curtz/tools/env"
	"github.com/sanctumlabs/curtz/tools/logger"
	"strconv"
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
	log := logger.NewLogger("vehicle-api")

	err := godotenv.Load()
	if err != nil {
		log.Warn("Error loading .env file. Using defaults")
	}

	environment := env.EnvOr(Env, "development")
	logLevel := env.EnvOr(EnvLogLevel, "debug")
	logJsonOutput := env.EnvOr(EnvLogJsonOutput, "true")
	port := env.EnvOr(EnvPort, "8080")
	host := env.EnvOr(EnvDatabaseHost, "localhost")
	database := env.EnvOr(EnvDatabase, "vehicles")
	databaseUser := env.EnvOr(EnvDatabaseUsername, "root")
	databasePass := env.EnvOr(EnvDatabasePassword, "root")
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

	// middlewares for the server
	corsMiddleware := middleware.NewCORSMiddleware(configuration.CorsHeaders)
	loggingMiddleware := middleware.NewLoggingMiddleware(configuration.Logging)
	recoveryMiddleware := middleware.NewRecoveryMiddleware()

	repository := repositories.NewRepository(configuration.Database)
	urlUseCase := url.NewUseCase(repository.GetUrlRepo())
	urlService := urlsvc.NewUrlService(urlUseCase)

	// setup routers
	routers := []router.Router{
		urlApi.NewUrlRouter(urlService),
		health.NewHealthRouter(),
	}

	// initialize routers
	srv.InitRouter(routers...)

	// use middlewares
	srv.UseMiddleware(loggingMiddleware)
	srv.UseMiddleware(corsMiddleware)
	srv.UseMiddleware(recoveryMiddleware)

	appServer := srv.CreateServer()

	// start & run the server
	err = appServer.Run(fmt.Sprintf(":%s", port))
	if err != nil {
		_, msg := fmt.Printf("Failed to start Server %s", err)
		log.Error(msg)
		panic(msg)
	}
}
