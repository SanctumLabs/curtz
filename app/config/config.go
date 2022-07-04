package config

type Config struct {
	Host        string
	Env         string
	Port        string
	Logging     LoggingConfig
	CorsHeaders string
	Version     string
	Database    DatabaseConfig
	Auth        AuthConfig
	Cache       CacheConfig
}
