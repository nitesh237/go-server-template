package log

import (
	"github.com/nitesh237/go-server-template/pkg/cfg"
	"go.uber.org/fx"
	"go.uber.org/fx/fxevent"
	"go.uber.org/zap"
)

var (
	FxEventZapLogger = fx.WithLogger(func(log *zap.Logger) fxevent.Logger {
		return &fxevent.ZapLogger{Logger: log}
	})

	FxZapModule = fx.Module("zap-logger",
		fx.Provide(
			NewZapLoggerImplProvider,
			NewLoggerProvierFromZapLoggerImpl,
			NewZapLoggerProviderFromZapLoggerImpl,
			NewUnwrappedZapLoggerProviderFromZapLoggerImpl,
		),
	)
)

type ZapLoggerProviderParams struct {
	fx.In

	Env         cfg.Environment `name:"Env"`
	Application *cfg.Application
}

type ZapLoggerProviderResult struct {
	fx.Out

	ZapLoggerImpl *ZapLoggerImpl
}

func NewZapLoggerImplProvider(p ZapLoggerProviderParams) (ZapLoggerProviderResult, error) {
	zapLogger, err := NewZapLogger(p.Env, p.Application.Logging)
	if err != nil {
		return ZapLoggerProviderResult{}, err
	}

	return ZapLoggerProviderResult{
		ZapLoggerImpl: zapLogger,
	}, nil
}

func NewLoggerProvierFromZapLoggerImpl(zapLoggerImpl *ZapLoggerImpl) Logger {
	return zapLoggerImpl
}

func NewZapLoggerProviderFromZapLoggerImpl(zapLoggerImpl *ZapLoggerImpl) ZapLogger {
	return zapLoggerImpl
}

func NewUnwrappedZapLoggerProviderFromZapLoggerImpl(zapLoggerImpl *ZapLoggerImpl) *zap.Logger {
	return zapLoggerImpl.logger
}
