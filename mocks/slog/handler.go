//revive:disable:package-comments
package slog

import (
	"context"
	"log/slog"
	"sync"
)

// CaptureHandler is a test helper that captures log records.
type CaptureHandler struct {
	records *[]slog.Record
	mu      sync.Mutex
}

// NewCaptureHandler creates a new CaptureHandler that appends records to the provided slice.
func NewCaptureHandler(records *[]slog.Record) *CaptureHandler {
	return &CaptureHandler{records: records}
}

// Enabled accepts every level, so nothing logged is lost to filtering.
func (h *CaptureHandler) Enabled(_ context.Context, _ slog.Level) bool {
	return true
}

// Handle appends r to the records slice the handler was built with.
func (h *CaptureHandler) Handle(_ context.Context, r slog.Record) error {
	h.mu.Lock()
	defer h.mu.Unlock()
	*h.records = append(*h.records, r)
	return nil
}

// WithAttrs answers with the same handler: attributes are not recorded
// separately from the records they would decorate.
func (h *CaptureHandler) WithAttrs(_ []slog.Attr) slog.Handler {
	return h
}

// WithGroup answers with the same handler: groups are not recorded.
func (h *CaptureHandler) WithGroup(_ string) slog.Handler {
	return h
}
