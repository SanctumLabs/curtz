package recoveryutils

import (
	"context"
	"fmt"
	"log/slog"
	"math"
	"math/rand/v2"
	"time"

	"github.com/sanctumlabs/curtz/app/pkg/errdefs"
)

type RetryConfig struct {
	MaxAttempts     int
	BaseDelay       time.Duration
	MaxDelay        time.Duration
	BackoffExponent float64
	JitterPercent   float64
}

var DefaultRetryConfig = RetryConfig{
	MaxAttempts:     3,
	BaseDelay:       100 * time.Millisecond,
	MaxDelay:        5 * time.Second,
	BackoffExponent: 2.0,
	JitterPercent:   0.1,
}

var CriticalRetryConfig = RetryConfig{
	MaxAttempts:     5,
	BaseDelay:       50 * time.Millisecond,
	MaxDelay:        3 * time.Second,
	BackoffExponent: 1.5,
	JitterPercent:   0.2,
}

// CalculateBackoffDelay calculates the delay for the next retry attempt
func (rc RetryConfig) CalculateBackoffDelay(attempt int) time.Duration {
	if attempt <= 0 {
		return rc.BaseDelay
	}

	// Exponential backoff
	delay := float64(rc.BaseDelay) * math.Pow(rc.BackoffExponent, float64(attempt-1))

	// Cap at max delay
	if delay > float64(rc.MaxDelay) {
		delay = float64(rc.MaxDelay)
	}

	// Add jitter to prevent thundering herd
	jitter := delay * rc.JitterPercent * (rand.Float64()*2 - 1) // +/- jitter percent
	delay += jitter

	// Ensure delay is not negative
	if delay < 0 {
		delay = float64(rc.BaseDelay)
	}

	return time.Duration(delay)
}

// RetryableOperation represents an operation that can be retried
type RetryableOperation[T any] func(ctx context.Context) (T, error)

// ExecuteWithRetry executes an operation with retry logic
func ExecuteWithRetry[T any](
	ctx context.Context,
	operation RetryableOperation[T],
	config RetryConfig,
	operationName string,
) (T, error) {
	var lastErr error
	var result T

	for attempt := 1; attempt <= config.MaxAttempts; attempt++ {
		select {
		case <-ctx.Done():
			return result, fmt.Errorf("operation cancelled during retry attempt %d: %w", attempt, ctx.Err())
		default:
		}

		result, lastErr = operation(ctx)
		if lastErr == nil {
			if attempt > 1 {
				slog.InfoContext(ctx, "Operation succeeded after retry",
					"operation", operationName,
					"attempt", attempt,
					"total_attempts", config.MaxAttempts)
			}
			return result, nil
		}

		if !errdefs.IsRetryableError(lastErr) {
			slog.ErrorContext(
				ctx,
				"Non-retryable error encountered",
				"operation", operationName,
				"attempt", attempt,
				"error", lastErr,
			)
			return result, fmt.Errorf("non-retryable error in %s: %w", operationName, lastErr)
		}

		if attempt < config.MaxAttempts {
			delay := config.CalculateBackoffDelay(attempt)
			slog.WarnContext(
				ctx,
				"Operation failed, retrying",
				"operation", operationName,
				"attempt", attempt,
				"total_attempts", config.MaxAttempts,
				"delay", delay,
				"error", lastErr,
			)

			select {
			case <-ctx.Done():
				return result, fmt.Errorf("operation cancelled during retry delay: %w", ctx.Err())
			case <-time.After(delay):
				// Continue to next attempt
			}
		}
	}

	slog.ErrorContext(
		ctx,
		"Operation failed after all retry attempts",
		"operation", operationName,
		"total_attempts", config.MaxAttempts,
		"final_error", lastErr,
	)

	return result, fmt.Errorf("operation %s failed after %d attempts: %w",
		operationName, config.MaxAttempts, lastErr)
}
