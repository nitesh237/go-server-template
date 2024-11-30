package goroutine

import "go.uber.org/fx"

var (
	FxSafeGoroutineWrapperModule = fx.Module("safe-goroutine-wrapper",
		fx.Provide(NewSafegoroutineWrapper),
		fx.Invoke(
			func(wrapper *safegoroutineWrapper) {
				defaultSafegoroutineWrapper = wrapper
			},
		),
	)
)
