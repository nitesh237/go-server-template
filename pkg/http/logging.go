package http

import "github.com/nitesh237/go-server-template/pkg/log"

type retryablehttpLeveledLogger struct {
	lg log.Logger
}

func (l *retryablehttpLeveledLogger) Error(msg string, keysAndValues ...interface{}) {
	l.lg.ErrorNoCtx(msg, keysAndValues...)
}

func (l *retryablehttpLeveledLogger) Info(msg string, keysAndValues ...interface{}) {
	l.lg.InfoNoCtx(msg, keysAndValues...)
}

func (l *retryablehttpLeveledLogger) Debug(msg string, keysAndValues ...interface{}) {
	l.lg.DebugNoCtx(msg, keysAndValues...)
}

func (l *retryablehttpLeveledLogger) Warn(msg string, keysAndValues ...interface{}) {
	l.lg.WarnNoCtx(msg, keysAndValues...)
}
