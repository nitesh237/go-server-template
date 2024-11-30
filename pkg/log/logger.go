package log

import (
	"context"

	"go.uber.org/zap"
)

type LoggerWithoutCtx interface {
	InfoNoCtx(msg string, a ...any)
	DebugNoCtx(msg string, a ...any)
	WarnNoCtx(msg string, a ...any)
	ErrorNoCtx(msg string, a ...any)
	PanicNoCtx(msg string, a ...any)
}

type Logger interface {
	LoggerWithoutCtx
	Info(ctx context.Context, msg string, a ...any)
	Debug(ctx context.Context, msg string, a ...any)
	Warn(ctx context.Context, msg string, a ...any)
	Error(ctx context.Context, msg string, a ...any)
	Panic(ctx context.Context, msg string, a ...any)
}

type ZapLogger interface {
	Unwrap() *zap.Logger
}
