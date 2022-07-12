package monitoring

import (
	"context"
	"time"

	"github.com/getsentry/sentry-go"
	"github.com/sanctumlabs/curtz/app/config"
	"github.com/sanctumlabs/curtz/app/tools/logger"
)

var log = logger.NewLogger("monitoring")

type Monitoring struct {
	cfg config.MonitoringConfig
}

// New creates a new monitoring utility
func New(config config.MonitoringConfig) (*Monitoring, error) {
	err := newSentry(config.Sentry)
	if err != nil {
		return nil, err
	}

	return &Monitoring{cfg: config}, nil
}

func ErrorHandler() {
	err := recover()

	if err != nil {
		log.Errorf("Recovering from err %s", err)
		sentry.CurrentHub().Recover(err)
		sentry.Flush(time.Second * 5)
	}
}

func RecoverWithContext(ctx context.Context) {
	sentry.RecoverWithContext(ctx)
}

func CaptureException(err error) {
	log.Errorf("Captured Exception err %s", err)
	sentry.CaptureException(err)
}

func CaptureMessage(message string) {
	log.Infof("Captured Message %s", message)
	sentry.CaptureMessage(message)
}
