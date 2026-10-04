package config

import (
	"errors"
	"net/url"
	"time"

	"github.com/sanctumlabs/curtz/app/pkg/infra/cache/redis"
	"github.com/sanctumlabs/curtz/app/pkg/infra/database/postgres"
)

const defaultMigrationsPath = "app/internal/adapters/postgres/migrations"

// ServerSettings configures the HTTP server.
type ServerSettings struct {
	Host    string
	Port    int
	Header  string
	Name    string
	Version string
	// BaseURL is the externally visible address, used to build links in emails.
	BaseURL string
}

// DatabaseSettings configures Postgres: the client settings plus the per-operation timeout the datastores use.
type DatabaseSettings struct {
	Postgres         postgres.PostgresDatabaseConfig
	OperationTimeout time.Duration
}

// MigrationSettings says where the SQL migrations live. Only the migrator reads it.
type MigrationSettings struct {
	Path string
}

// App is everything the API process reads from its environment.
type App struct {
	Environment     string
	Server          ServerSettings
	Database        DatabaseSettings
	Redis           redis.RedisClientConfig
	Auth            AuthConfig
	ShutdownTimeout time.Duration
}

// Load reads and validates the whole application configuration. Every problem is reported, not just the first.
func Load(lookup Lookup) (App, error) {
	r := newReader(lookup)
	app := App{
		Environment:     r.str("ENVIRONMENT", environmentDevelopment),
		ShutdownTimeout: r.units("SHUTDOWN_TIMEOUT", 15, time.Second),
	}
	if app.ShutdownTimeout <= 0 {
		r.fail("SHUTDOWN_TIMEOUT must be greater than zero")
	}

	errs := []error{r.err()}
	var err error
	app.Server, err = LoadServer(lookup)
	errs = append(errs, err)
	app.Database, err = LoadDatabase(lookup)
	errs = append(errs, err)
	app.Redis, err = LoadRedis(lookup)
	errs = append(errs, err)
	app.Auth, err = LoadAuth(lookup)
	errs = append(errs, err)

	return app, errors.Join(errs...)
}

// LoadServer reads the HTTP server settings.
func LoadServer(lookup Lookup) (ServerSettings, error) {
	r := newReader(lookup)
	settings := ServerSettings{
		Host:    r.str("SERVER_HOST", "0.0.0.0"),
		Port:    r.integer("HTTP_PORT", 8085),
		Header:  r.str("SERVER_HEADER", "Curtz"),
		Name:    r.str("SERVER_NAME", "Curtz"),
		Version: r.str("SERVER_VERSION", "1.0.0"),
		BaseURL: r.str("APP_BASE_URL", "http://localhost:8085"),
	}
	if settings.Port < 1 || settings.Port > 65535 {
		r.fail("HTTP_PORT must be a port between 1 and 65535")
	}
	return settings, r.err()
}

// LoadDatabase reads the Postgres settings. The migrator uses only this loader, so it needs no AUTH_SECRET.
func LoadDatabase(lookup Lookup) (DatabaseSettings, error) {
	r := newReader(lookup)
	pg := postgres.PostgresDatabaseConfig{
		Host:            r.str("DATABASE_HOST", "localhost"),
		Port:            r.str("DATABASE_PORT", "5432"),
		Name:            r.str("DATABASE_NAME", "curtzdb"),
		Username:        r.str("DATABASE_USERNAME", "curtz-user"),
		Password:        r.str("DATABASE_PASSWORD", devDatabasePassword),
		Url:             r.str("DATABASE_URL", ""),
		SslMode:         r.str("DATABASE_SSL_MODE", "disable"),
		MaxConns:        int32(r.integer("DATABASE_MAX_CONNS", 30)),
		MinConns:        int32(r.integer("DATABASE_MIN_CONNS", 5)),
		MaxConnLifetime: r.units("DATABASE_MAX_CONN_LIFETIME", 1, time.Hour),
		MaxConnIdleTime: r.units("DATABASE_MAX_CONN_IDLE_TIME", 30, time.Minute),
		ConnTimeout:     r.units("DATABASE_CONN_TIMEOUT", 30, time.Second),
		QueryTimeout:    r.units("DATABASE_QUERY_TIMEOUT", 10, time.Second),
	}
	settings := DatabaseSettings{
		Postgres:         pg,
		OperationTimeout: r.units("DATABASE_OPERATION_TIMEOUT", 30, time.Second),
	}

	r.port("DATABASE_PORT", pg.Port)
	if pg.MaxConns < 1 {
		r.fail("DATABASE_MAX_CONNS must be at least 1")
	}
	if pg.MinConns > pg.MaxConns {
		r.fail("DATABASE_MIN_CONNS must not exceed DATABASE_MAX_CONNS")
	}

	password := pg.Password
	if pg.Url != "" {
		parsed, err := url.Parse(pg.Url)
		if err != nil || parsed.Scheme == "" || parsed.Host == "" {
			// The parse error text can contain the URL, and with it the password, so it is not repeated.
			r.fail("DATABASE_URL is not a valid database URL")
		} else if parsed.User != nil {
			password, _ = parsed.User.Password()
		} else {
			password = ""
		}
	}
	if r.enforceSecrets() && password == devDatabasePassword {
		if pg.Url != "" {
			r.fail("DATABASE_URL must not use the development password when ENVIRONMENT is not development or test")
		} else {
			r.fail("DATABASE_PASSWORD must be set to a non-default value when ENVIRONMENT is not development or test")
		}
	}

	return settings, r.err()
}

// LoadRedis reads the Redis settings. REDIS_ADDRESS is a comma-separated host:port list: one entry gives a plain
// client, several give a cluster client.
func LoadRedis(lookup Lookup) (redis.RedisClientConfig, error) {
	r := newReader(lookup)
	cfg := redis.RedisClientConfig{
		Address:  splitList(r.str("REDIS_ADDRESS", "localhost:7001")),
		Username: r.str("REDIS_USERNAME", "curtz-svc"),
		Password: r.str("REDIS_PASSWORD", devRedisPassword),
		Database: r.integer("REDIS_DATABASE", 0),
	}

	if len(cfg.Address) == 0 {
		r.fail("REDIS_ADDRESS must list at least one host:port")
	}
	for _, entry := range cfg.Address {
		if !validAddress(entry) {
			r.fail("REDIS_ADDRESS entry %q must be host:port with a port between 1 and 65535", entry)
		}
	}
	if cfg.Database != 0 {
		r.fail("REDIS_DATABASE must be 0 (the Redis Cluster serves database 0 only)")
	}
	if r.enforceSecrets() && cfg.Password == devRedisPassword {
		r.fail("REDIS_PASSWORD must be set to a non-default value when ENVIRONMENT is not development or test")
	}

	return cfg, r.err()
}

// LoadAuth reads the JWT settings.
func LoadAuth(lookup Lookup) (AuthConfig, error) {
	r := newReader(lookup)
	cfg := AuthConfig{Jwt: Jwt{
		Secret:             r.str("AUTH_SECRET", devAuthSecret),
		Issuer:             r.str("AUTH_ISSUER", "curtz"),
		ExpireDelta:        r.integer("AUTH_EXPIRE_DELTA", 15),
		RefreshExpireDelta: r.integer("AUTH_REFRESH_EXPIRE_DELTA", 24),
	}}

	if cfg.ExpireDelta < 1 {
		r.fail("AUTH_EXPIRE_DELTA must be at least 1")
	}
	if cfg.RefreshExpireDelta < 1 {
		r.fail("AUTH_REFRESH_EXPIRE_DELTA must be at least 1")
	}
	if r.enforceSecrets() && cfg.Secret == devAuthSecret {
		r.fail("AUTH_SECRET must be set to a non-default value when ENVIRONMENT is not development or test")
	}

	return cfg, r.err()
}

// LoadMigrations reads where the SQL migrations live.
func LoadMigrations(lookup Lookup) (MigrationSettings, error) {
	r := newReader(lookup)
	return MigrationSettings{Path: r.str("MIGRATIONS_PATH", defaultMigrationsPath)}, r.err()
}
