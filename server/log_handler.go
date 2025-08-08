package server

import (
	"context"
	"log/slog"

	"cosmossdk.io/log"
)

// CustomSlogHandler bridges Geth's slog logs to the existing Cosmos SDK logger.
type CustomSlogHandler struct {
	logger log.Logger
}

// Handle processes slog records and forwards them to your Cosmos SDK logger.
func (h *CustomSlogHandler) Handle(_ context.Context, r slog.Record) error {
	attrs := []interface{}{}
	r.Attrs(func(attr slog.Attr) bool {
		attrs = append(attrs, attr.Key, attr.Value.Any())
		return true
	})

	// Map slog levels to Cosmos SDK logger
	switch r.Level {
	case slog.LevelDebug:
		h.logger.Debug(r.Message, attrs...)
	case slog.LevelInfo:
		h.logger.Info(r.Message, attrs...)
	case slog.LevelWarn:
		h.logger.Warn(r.Message, attrs...)
	case slog.LevelError:
		h.logger.Error(r.Message, attrs...)
	default:
		h.logger.Info(r.Message, attrs...)
	}

	return nil
}

// Enabled determines if the handler should log a given level.
func (h *CustomSlogHandler) Enabled(_ context.Context, _ slog.Level) bool {
	return true
}

// WithAttrs allows adding additional attributes.
func (h *CustomSlogHandler) WithAttrs(_ []slog.Attr) slog.Handler {
	return h
}

// WithGroup is required to implement slog.Handler (not used).
func (h *CustomSlogHandler) WithGroup(_ string) slog.Handler {
	return h
}
