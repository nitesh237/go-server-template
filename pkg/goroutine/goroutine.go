package goroutine

import (
	"fmt"

	"github.com/nitesh237/go-server-template/pkg/log"
	"go.uber.org/zap"
)

var (
	defaultSafegoroutineWrapper *safegoroutineWrapper
)

func NewSafegoroutineWrapper(lg log.Logger) *safegoroutineWrapper {
	return &safegoroutineWrapper{
		lg: lg,
	}
}

type safegoroutineWrapper struct {
	lg log.Logger
}

/*
 * `Go` provides a safe way to execute a function asynchronously, recovering if the panic might occur.
 */
func (g *safegoroutineWrapper) Go(fn func()) {
	go func(lg log.Logger) {
		defer recoverPanic(lg)
		fn()
	}(g.lg)
}

/*
 * Write the error to console when a goroutine of a task panicking.
 */
func recoverPanic(lg log.Logger) {
	if r := recover(); r != nil {
		err, ok := r.(error)
		if !ok {
			err = fmt.Errorf("%v", r)
		}
		lg.ErrorNoCtx("goroutine panic", zap.Error(err))
	}
}

// `Go` provides a safe way to execute a function asynchronously, recovering if the panic might occur.
func Go(fn func()) {
	defaultSafegoroutineWrapper.Go(fn)
}
