package config

type Sentry struct {
	DSN              string
	TracesSampleRate float64
	Environment      string
	Enabled          bool
}

type MonitoringConfig struct {
	Sentry
}
