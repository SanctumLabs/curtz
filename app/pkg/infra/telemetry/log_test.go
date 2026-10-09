package telemetry

import (
	"bytes"
	"context"
	"encoding/json"
	"log/slog"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/otel/trace"
)

const (
	testTraceID = "4bf92f3577b34da6a3ce929d0e0e4736"
	testSpanID  = "00f067aa0ba902b7"
)

// spanContext returns a context that carries a span, sampled or not.
func spanContext(t *testing.T, flags trace.TraceFlags) context.Context {
	t.Helper()
	traceID, err := trace.TraceIDFromHex(testTraceID)
	require.NoError(t, err)
	spanID, err := trace.SpanIDFromHex(testSpanID)
	require.NoError(t, err)
	return trace.ContextWithSpanContext(context.Background(), trace.NewSpanContext(trace.SpanContextConfig{
		TraceID: traceID, SpanID: spanID, TraceFlags: flags,
	}))
}

func logLine(t *testing.T, buf *bytes.Buffer) map[string]any {
	t.Helper()
	var line map[string]any
	require.NoError(t, json.Unmarshal(buf.Bytes(), &line), buf.String())
	return line
}

func TestNewLogger_JSONCarriesTheKeysTheELKPipelineMaps(t *testing.T) {
	var buf bytes.Buffer
	logger := NewLogger(&buf, "json", slog.LevelInfo, "curtz")

	logger.InfoContext(spanContext(t, trace.FlagsSampled), "hello", "user", "u-1")

	line := logLine(t, &buf)
	assert.Equal(t, "hello", line["msg"])
	assert.Equal(t, "INFO", line["level"])
	assert.NotEmpty(t, line["time"])
	assert.Equal(t, "curtz", line["service"])
	assert.Equal(t, "u-1", line["user"])
	assert.Equal(t, testTraceID, line["trace_id"])
	assert.Equal(t, testSpanID, line["span_id"])
}

func TestNewLogger_OmitsTheIDsWithoutASpan(t *testing.T) {
	var buf bytes.Buffer
	logger := NewLogger(&buf, "json", slog.LevelInfo, "curtz")

	logger.InfoContext(context.Background(), "no span")
	logger.Info("no context at all")

	for _, raw := range strings.Split(strings.TrimSpace(buf.String()), "\n") {
		assert.NotContains(t, raw, "trace_id")
		assert.NotContains(t, raw, "span_id")
		assert.Contains(t, raw, `"service":"curtz"`)
	}
}

// A span that was sampled out still has IDs, and the log line is still worth finding by them.
func TestNewLogger_AddsTheIDsOfAnUnsampledSpanToo(t *testing.T) {
	var buf bytes.Buffer
	logger := NewLogger(&buf, "json", slog.LevelInfo, "curtz")

	logger.InfoContext(spanContext(t, 0), "sampled out")

	assert.Equal(t, testTraceID, logLine(t, &buf)["trace_id"])
}

// slog.With wraps the handler through WithAttrs; a wrapper that forgot to wrap its result would silently stop adding IDs.
func TestNewLogger_DerivedLoggersKeepAddingTheIDs(t *testing.T) {
	var buf bytes.Buffer
	logger := NewLogger(&buf, "json", slog.LevelInfo, "curtz")

	logger.With("component", "x").WithGroup("g").InfoContext(spanContext(t, trace.FlagsSampled), "derived", "k", "v")

	line := logLine(t, &buf)
	assert.Equal(t, "x", line["component"])
	group, ok := line["g"].(map[string]any)
	require.True(t, ok, "WithGroup must still group the call's attributes")
	assert.Equal(t, "v", group["k"])
	// Handle runs inside the open group, so the IDs land in it; what matters is that they are not lost.
	assert.Equal(t, testTraceID, group["trace_id"])
}

func TestNewLogger_HonoursTheLevel(t *testing.T) {
	var buf bytes.Buffer
	logger := NewLogger(&buf, "json", slog.LevelWarn, "curtz")

	logger.Info("dropped")
	logger.Debug("dropped")
	assert.Empty(t, buf.String())

	logger.Warn("kept")
	assert.Contains(t, buf.String(), "kept")
}

func TestNewLogger_TextFormatIsForTerminals(t *testing.T) {
	var buf bytes.Buffer
	logger := NewLogger(&buf, "text", slog.LevelInfo, "curtz")

	logger.InfoContext(spanContext(t, trace.FlagsSampled), "hello")

	out := buf.String()
	assert.False(t, strings.HasPrefix(out, "{"), "text format must not be JSON: %s", out)
	assert.Contains(t, out, "msg=hello")
	assert.Contains(t, out, "service=curtz")
	assert.Contains(t, out, "trace_id="+testTraceID)
}

func TestServiceName(t *testing.T) {
	t.Setenv("OTEL_SERVICE_NAME", "")
	assert.Equal(t, "curtz", ServiceName(DefaultServiceName), "an empty value counts as unset")
	assert.Equal(t, "curtz-worker", ServiceName("curtz-worker"))

	t.Setenv("OTEL_SERVICE_NAME", "  curtz-api ")
	assert.Equal(t, "curtz-api", ServiceName("curtz-worker"), "the variable wins over the fallback")
}
