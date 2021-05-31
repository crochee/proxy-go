package logger

import "context"

type loggerKey struct{}

// Context Adds fields.
func Context(ctx context.Context, log Builder) context.Context {
	return context.WithValue(ctx, loggerKey{}, log)
}

// FromContext Gets the logger from context.
func FromContext(ctx context.Context) Builder {
	l, ok := ctx.Value(loggerKey{}).(Builder)
	if !ok {
		return NoLogger{}
	}
	return l
}
