package config

type Config struct {
	Env         string
	Port        string
	Logging     LoggingConfig
	CorsHeaders string
	Version     string
	Database    DatabaseConfig
}
