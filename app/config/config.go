package config

type Config struct {
	Host        string
	Env         string
	DocsEnabled bool
	Port        string
	Logging     LoggingConfig
	CorsHeaders string
	Version     string
	Database    DatabaseConfig
	Auth        AuthConfig
	Cache       CacheConfig
	Monitoring  MonitoringConfig
}
