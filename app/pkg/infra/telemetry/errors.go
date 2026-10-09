package telemetry

import (
	"log/slog"
	"sync"
	"time"
)

// maxTrackedErrors bounds the memory of the handler: when more distinct errors than this are seen within one interval the
// oldest records are forgotten, which at worst repeats a line.
const maxTrackedErrors = 64

// errorHandler is the OpenTelemetry error handler. The SDK reports every failed export to it, which for an unreachable
// collector means one error every few seconds for as long as the collector is away. It logs each distinct error at
// most once per interval, so the log carries one line a minute instead.
type errorHandler struct {
	mu       sync.Mutex
	interval time.Duration
	now      func() time.Time
	lastSeen map[string]time.Time
}

func newErrorHandler(interval time.Duration, now func() time.Time) *errorHandler {
	return &errorHandler{interval: interval, now: now, lastSeen: map[string]time.Time{}}
}

// Handle implements otel.ErrorHandler.
func (h *errorHandler) Handle(err error) {
	if err == nil {
		return
	}
	if !h.admit(err.Error()) {
		return
	}
	slog.Warn("telemetry export failed", "error", err)
}

// admit reports whether message has not been logged within the interval, and records that it is being logged now.
func (h *errorHandler) admit(message string) bool {
	h.mu.Lock()
	defer h.mu.Unlock()

	now := h.now()
	if last, seen := h.lastSeen[message]; seen && now.Sub(last) < h.interval {
		return false
	}
	if len(h.lastSeen) >= maxTrackedErrors {
		for key, last := range h.lastSeen {
			if now.Sub(last) >= h.interval {
				delete(h.lastSeen, key)
			}
		}
		if len(h.lastSeen) >= maxTrackedErrors {
			clear(h.lastSeen)
		}
	}
	h.lastSeen[message] = now
	return true
}
