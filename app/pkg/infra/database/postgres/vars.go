package postgres

import "time"

const (
	defaultConnAttempts               = 3
	defaultConnTimeout  time.Duration = time.Second

	OperationTypeDefault  = "default"
	OperationTypeCritical = "critical"
	OperationTypeRead     = "read"
	OperationTypeBulk     = "bulk"
)
