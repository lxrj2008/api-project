package utils

import "context"

// ctxKey is a private type preventing key collisions.
type ctxKey string

const (
	ctxKeyTraceID ctxKey = "trace_id"
	ctxKeyUserID  ctxKey = "user_id"
	GinKeyTraceID        = "traceId"
	GinKeyUserID         = "userId"
)

// WithTraceID returns a new context carrying trace ID.
func WithTraceID(ctx context.Context, traceID string) context.Context {
	return context.WithValue(ctx, ctxKeyTraceID, traceID)
}

// TraceIDFromContext extracts trace ID if present.
func TraceIDFromContext(ctx context.Context) string {
	if v, ok := ctx.Value(ctxKeyTraceID).(string); ok {
		return v
	}
	return ""
}

// WithUserID stores authenticated user ID.
func WithUserID(ctx context.Context, userID string) context.Context {
	return context.WithValue(ctx, ctxKeyUserID, userID)
}

// UserIDFromContext fetches user identifier.
func UserIDFromContext(ctx context.Context) string {
	if v, ok := ctx.Value(ctxKeyUserID).(string); ok {
		return v
	}
	return ""
}
