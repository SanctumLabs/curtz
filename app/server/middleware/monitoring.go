package middleware

import (
	sentrygin "github.com/getsentry/sentry-go/gin"
	"github.com/gin-gonic/gin"
	"github.com/sanctumlabs/curtz/app/config"
	"github.com/sanctumlabs/curtz/app/tools/logger"
	"github.com/sanctumlabs/curtz/app/tools/monitoring"
)

func NewMonitoringMiddleware(config config.MonitoringConfig) Middleware {
	log := logger.NewLogger("monitoring-log")
	return func(context *gin.Context) {
		if config.Sentry.Enabled && config.Sentry.DSN != "" {
			if err := monitoring.NewSentry(config.Sentry, context); err != nil {
				log.Fatalf("Failed to configure Sentry: %s", err)
			}

			sentrygin.New(sentrygin.Options{
				Repanic: true,
			})

			if hub := sentrygin.GetHubFromContext(context); hub != nil {
				hub.Scope().SetTag("Testing", "Testing integration")
			}

			context.Next()
		}
	}
}
