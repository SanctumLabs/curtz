package monitoring

import (
	"time"

	"github.com/getsentry/sentry-go"
	"github.com/sanctumlabs/curtz/app/config"
)

func newSentry(config config.Sentry) error {
	if config.Enabled && config.DSN != "" {
		debug := config.Environment == "development"

		if err := sentry.Init(sentry.ClientOptions{
			Dsn:              config.DSN,
			Environment:      config.Environment,
			TracesSampleRate: config.TracesSampleRate,
			Debug:            debug,
			AttachStacktrace: true,
			BeforeSend: func(event *sentry.Event, hint *sentry.EventHint) *sentry.Event {
				event.User.Email = ""
				return event
			},
		}); err != nil {
			log.Fatalf("SentryInit: %s", err)
			return err
		}

		// Flush buffered events before the program terminates.
		// Set the timeout to the maximum duration the program can afford to wait.
		defer sentry.Flush(2 * time.Second)
		return nil
	}
	return nil
}
