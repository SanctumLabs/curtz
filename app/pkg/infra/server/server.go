package server

import (
	"carduka/bidsvc/pkg/infra/logger"
	"carduka/bidsvc/pkg/infra/server/middleware"
	"carduka/bidsvc/pkg/infra/server/router"
	"fmt"

	"github.com/bytedance/sonic"
	"github.com/gofiber/contrib/swagger"
	"github.com/gofiber/fiber/v2"
)

type Server struct {
	app *fiber.App
	cfg ServerConfig
	log logger.Logger
}

// NewServer creates a new server
func NewServer(cfg ServerConfig) *Server {
	appLogger := logger.New(nil)

	app := fiber.New(fiber.Config{
		ServerHeader: cfg.Header,
		AppName:      cfg.AppName,
		// Using a custom encoder and decoder to marshal and unmarshal JSON
		// ref: https://github.com/bytedance/sonic
		JSONEncoder: sonic.Marshal,
		JSONDecoder: sonic.Unmarshal,
	})

	// middleware
	app.Use(middleware.RequestIdMiddleware())
	app.Use(middleware.LoggerMiddleware())
	app.Use(middleware.CORSMiddleware())

	// TODO: File path is relative to the binary created. This needs to be set accordingly depending on whether running
	// the binary from the root of the project or from the build directory

	// v1 of the documentation
	app.Use(middleware.SwaggerMiddleware(
		swagger.Config{
			BasePath: "/",
			Path:     "docs/v1",
			FilePath: "./api/openapi-spec/bids_service_v1.swagger.json",
			Title:    "Bids Service V1 API Docs",
		},
	))

	app.Use(middleware.SwaggerMiddleware(
		swagger.Config{
			BasePath: "/",
			Path:     "docs/monitoring",
			FilePath: "./api/openapi-spec/monitoring.swagger.json",
			Title:    "Bids Service Monitoring API Docs",
		},
	))

	app.Use(middleware.HelmetMiddleware())
	app.Use(middleware.IdempotencyMiddleware())

	app.Get("/metrics", middleware.MonitoringMiddleware())

	app.Use(middleware.RecoverMiddleware())

	return &Server{
		app: app,
		cfg: cfg,
		log: appLogger,
	}
}

func (srv *Server) Listen() error {
	srv.log.Infow("Listening on port", "port", srv.cfg.Port)
	return srv.app.Listen(fmt.Sprintf(":%d", srv.cfg.Port))
}

// Shutdown shutdowns the server
func (srv *Server) Shutdown() error {
	srv.log.Infow("shutting down server", "port", srv.cfg.Port)
	return srv.app.Shutdown()
}

// RegisterHandlers registers all the handlers for the user v1 endpoint
func (srv *Server) RegisterHandlers(router []router.Router) {
	for _, r := range router {
		routes := r.Routes()
		for _, route := range routes {
			srv.app.Add(route.Method(), route.Path(), route.Handler())
		}
	}
}
