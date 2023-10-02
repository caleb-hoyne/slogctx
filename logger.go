package slogctx

import (
	"context"
	"log/slog"
)

type contextKey struct{}

type Handler struct {
	slog.Handler
}

func (h *Handler) Handle(ctx context.Context, r slog.Record) error {
	r.AddAttrs(h.observe(ctx)...)
	return h.Handler.Handle(ctx, r)
}

func (h *Handler) observe(ctx context.Context) []slog.Attr {
	if ctx == nil {
		return nil
	}
	if v := ctx.Value(contextKey{}); v != nil {
		if attrs, ok := v.([]slog.Attr); ok {
			return attrs
		}
	}
	return nil
}

func AddValues(ctx context.Context, attrs ...slog.Attr) context.Context {
	if len(attrs) == 0 {
		return ctx
	}
	if ctx == nil {
		ctx = context.Background()
	}
	if v := ctx.Value(contextKey{}); v != nil {
		if old, ok := v.([]slog.Attr); ok {
			attrs = append(old, attrs...)
		}
	}
	return context.WithValue(ctx, contextKey{}, attrs)
}
