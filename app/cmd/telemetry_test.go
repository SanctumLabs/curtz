package main

import (
	"bytes"
	"context"
	"errors"
	"log/slog"
	"testing"

	"github.com/sanctumlabs/curtz/app/config"
	"github.com/sanctumlabs/curtz/app/pkg"
	"github.com/sanctumlabs/curtz/app/pkg/infra/telemetry"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type setupFunc = func(context.Context, telemetry.Options) (func(context.Context) error, error)

func stubSetup(t *testing.T, stub setupFunc) {
	t.Helper()
	previous := setupTelemetry
	setupTelemetry = stub
	t.Cleanup(func() { setupTelemetry = previous })
}

func captureLog(t *testing.T) *bytes.Buffer {
	t.Helper()
	var buf bytes.Buffer
	previous := slog.Default()
	slog.SetDefault(slog.New(slog.NewTextHandler(&buf, nil)))
	t.Cleanup(func() { slog.SetDefault(previous) })
	return &buf
}

func TestStartTelemetry_PassesTheVersionAndTheEnvironment(t *testing.T) {
	var got telemetry.Options
	stubSetup(t, func(_ context.Context, opts telemetry.Options) (func(context.Context) error, error) {
		got = opts
		return func(context.Context) error { return nil }, nil
	})

	startTelemetry(context.Background(), config.App{Environment: "staging"})()

	assert.Equal(t, telemetry.Options{ServiceVersion: pkg.Version, Environment: "staging"}, got)
}

// The run context is cancelled by the SIGTERM that starts the shutdown, so a flush that reused it would give up at once.
func TestStartTelemetry_TheFlushGetsItsOwnDeadlineAfterTheRunContextIsCancelled(t *testing.T) {
	var flushErr error
	var hasDeadline bool
	stubSetup(t, func(context.Context, telemetry.Options) (func(context.Context) error, error) {
		return func(ctx context.Context) error {
			flushErr = ctx.Err()
			_, hasDeadline = ctx.Deadline()
			return nil
		}, nil
	})
	ctx, cancel := context.WithCancel(context.Background())
	flush := startTelemetry(ctx, config.App{})
	cancel() // the SIGTERM

	flush()

	require.NoError(t, flushErr, "the flush context must not be the cancelled run context")
	assert.True(t, hasDeadline, "and it must be bounded")
}

// A telemetry problem must never stop the API.
func TestStartTelemetry_ASetupFailureDisablesTelemetryInsteadOfStoppingTheAPI(t *testing.T) {
	buf := captureLog(t)
	stubSetup(t, func(context.Context, telemetry.Options) (func(context.Context) error, error) {
		return nil, errors.New("bad exporter configuration")
	})

	flush := startTelemetry(context.Background(), config.App{})

	require.NotNil(t, flush)
	assert.NotPanics(t, flush)
	assert.Contains(t, buf.String(), "telemetry is disabled")
	assert.Contains(t, buf.String(), "bad exporter configuration")
}

func TestStartTelemetry_AFailingFlushIsLoggedNotFatal(t *testing.T) {
	buf := captureLog(t)
	stubSetup(t, func(context.Context, telemetry.Options) (func(context.Context) error, error) {
		return func(context.Context) error { return errors.New("collector unreachable") }, nil
	})

	assert.NotPanics(t, startTelemetry(context.Background(), config.App{}))

	assert.Contains(t, buf.String(), "flushing telemetry")
	assert.Contains(t, buf.String(), "collector unreachable")
}

// run flushes right after the server drains and again from a defer that covers its early returns.
func TestStartTelemetry_TheFlushRunsOnlyOnce(t *testing.T) {
	var shutdowns int
	stubSetup(t, func(context.Context, telemetry.Options) (func(context.Context) error, error) {
		return func(context.Context) error {
			shutdowns++
			return nil
		}, nil
	})

	flush := startTelemetry(context.Background(), config.App{})
	flush()
	flush()

	assert.Equal(t, 1, shutdowns)
}
