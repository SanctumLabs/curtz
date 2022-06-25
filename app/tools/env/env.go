package env

import (
	"os"
)

// EnvOr gets an environment variable by the given key, if not found def is returned
func EnvOr(key, def string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return def
}
