package database

import (
	"time"

	recoveryutils "github.com/sanctumlabs/curtz/app/pkg/utils/recover"
)

// Config is the base configuration for handling database connections. This can be used by database implementations to configure how
// database transactions are created and handled or how retry mechanism can be handled.
type Config struct {
	// OperationTimeout is how long a database timeout is
	OperationTimeout time.Duration

	// RetryConfig retry configuration for the database operation
	RetryConfig recoveryutils.RetryConfig
}
