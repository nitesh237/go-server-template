package log

import (
	"context"

	"github.com/nitesh237/go-server-template/pkg/cfg"
	"github.com/nitesh237/go-server-template/pkg/errors"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

type ZapLoggerImpl struct {
	logger *zap.Logger
}

// NewZapLogger initialised logger implementation using zap
func NewZapLogger(env cfg.Environment, conf *cfg.Logging) (*ZapLoggerImpl, error) {
	var (
		logger *zap.Logger
		err    error
	)
	if conf.EnableLoggingToFile {
		logger, _, err = newZapFileLogger(env, conf)
	} else {
		logger, _, err = newZapLogger(env)
	}

	if err != nil {
		return nil, errors.Wrap(err, "failed to initialise logger")
	}

	return &ZapLoggerImpl{logger: logger}, nil
}

func (l *ZapLoggerImpl) InfoNoCtx(msg string, a ...any) {
	l.logger.Info(msg, _getZapFieldsFromGenerics(a...)...)
}

func (l *ZapLoggerImpl) DebugNoCtx(msg string, a ...any) {
	l.logger.Debug(msg, _getZapFieldsFromGenerics(a...)...)
}

func (l *ZapLoggerImpl) WarnNoCtx(msg string, a ...any) {
	l.logger.Warn(msg, _getZapFieldsFromGenerics(a...)...)
}

func (l *ZapLoggerImpl) ErrorNoCtx(msg string, a ...any) {
	l.logger.Error(msg, _getZapFieldsFromGenerics(a...)...)
}

func (l *ZapLoggerImpl) PanicNoCtx(msg string, a ...any) {
	l.logger.Panic(msg, _getZapFieldsFromGenerics(a...)...)
}

func (l *ZapLoggerImpl) Info(ctx context.Context, msg string, a ...any) {
	l.InfoNoCtx(msg, a...)
}

func (l *ZapLoggerImpl) Debug(ctx context.Context, msg string, a ...any) {
	l.DebugNoCtx(msg, a...)
}

func (l *ZapLoggerImpl) Warn(ctx context.Context, msg string, a ...any) {
	l.WarnNoCtx(msg, a...)
}

func (l *ZapLoggerImpl) Error(ctx context.Context, msg string, a ...any) {
	l.ErrorNoCtx(msg, a...)
}

func (l *ZapLoggerImpl) Panic(ctx context.Context, msg string, a ...any) {
	l.PanicNoCtx(msg, a...)
}

func (l *ZapLoggerImpl) Log(keyvals ...interface{}) error {
	l.logger.Sugar().Info(keyvals...)
	return nil
}

func (l *ZapLoggerImpl) Unwrap() *zap.Logger {
	return l.logger
}

func newZapLogger(env cfg.Environment) (*zap.Logger, zap.AtomicLevel, error) {
	switch env {
	case cfg.Dev, cfg.Test, cfg.Docker:
		return _newDevelopment()
	case cfg.QA:
		return _newProductionWithDebug()
	default:
		return _newProduction()
	}
}

func newZapFileLogger(env cfg.Environment, loggerConfig *cfg.Logging) (*zap.Logger, zap.AtomicLevel, error) {
	switch env {
	case cfg.Dev, cfg.Test, cfg.Docker:
		return _newDevelopment()
	case cfg.QA:
		return _newProductionWithDebugToFile(loggerConfig)
	default:
		return _newProductionToFile(loggerConfig)
	}
}

// _newDevelopment builds a custom logger for development environment which skips this wrapper's line number,
// file name etc while logging
func _newDevelopment(options ...zap.Option) (*zap.Logger, zap.AtomicLevel, error) {
	config := zap.NewDevelopmentConfig()
	options = append(options, zap.AddCallerSkip(1))

	logger, err := config.Build(options...)
	if err != nil {
		return nil, zap.AtomicLevel{}, err
	}

	return logger, config.Level, nil
}

// _newProductionWithDebug builds a custom logger for qa environment that writes DebugLevel and above logs
// to standard error as JSON.
func _newProductionWithDebug(options ...zap.Option) (*zap.Logger, zap.AtomicLevel, error) {
	config := zap.NewProductionConfig()
	config.Development = true
	config.Level = zap.NewAtomicLevelAt(zap.DebugLevel)
	config.EncoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	options = append(options, zap.AddCallerSkip(1))

	logger, err := config.Build(options...)
	if err != nil {
		return nil, zap.AtomicLevel{}, err
	}

	return logger, config.Level, nil
}

// _newProduction builds a custom logger for production environment which skips this wrapper's line number,
// file name etc while logging
func _newProduction(options ...zap.Option) (*zap.Logger, zap.AtomicLevel, error) {
	config := zap.NewProductionConfig()
	config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	options = append(options, zap.AddCallerSkip(1))

	logger, err := config.Build(options...)
	if err != nil {
		return nil, zap.AtomicLevel{}, err
	}

	return logger, config.Level, nil
}

// _newProductionWithDebugToFile builds a custom logger for staging/demo environment that writes DebugLevel and above logs
// to standard error as JSON. It prints the log to the input file
func _newProductionWithDebugToFile(loggerConfig *cfg.Logging, options ...zap.Option) (*zap.Logger, zap.AtomicLevel, error) {
	w := zapcore.AddSync(&lumberjack.Logger{
		Filename:   loggerConfig.LogPath,
		MaxSize:    loggerConfig.MaxSizeInMBs,
		MaxBackups: loggerConfig.MaxBackups,
	})
	config := zap.NewProductionEncoderConfig()
	config.EncodeLevel = zapcore.CapitalLevelEncoder
	config.EncodeTime = zapcore.ISO8601TimeEncoder
	atomicLogLevel := zap.NewAtomicLevelAt(zap.DebugLevel)
	core := zapcore.NewCore(zapcore.NewJSONEncoder(config), w, atomicLogLevel)
	options = append(options, zap.AddCaller(), zap.AddStacktrace(zap.ErrorLevel), zap.AddCallerSkip(1), zap.ErrorOutput(w))
	logger := zap.New(core, options...)
	return logger, atomicLogLevel, nil
}

// _newProductionToFile builds a custom logger for production environment which skips this wrapper's line number,
// file name etc while logging. It prints the log to the input file
func _newProductionToFile(loggerConfig *cfg.Logging, options ...zap.Option) (*zap.Logger, zap.AtomicLevel, error) {
	w := zapcore.AddSync(&lumberjack.Logger{
		Filename:   loggerConfig.LogPath,
		MaxSize:    loggerConfig.MaxSizeInMBs,
		MaxBackups: loggerConfig.MaxBackups,
	})

	config := zap.NewProductionEncoderConfig()
	config.EncodeTime = zapcore.ISO8601TimeEncoder
	atomicLogLevel := zap.NewAtomicLevelAt(zap.InfoLevel)
	core := zapcore.NewCore(zapcore.NewJSONEncoder(config), w, atomicLogLevel)
	options = append(options, zap.AddCaller(), zap.AddStacktrace(zap.ErrorLevel), zap.AddCallerSkip(1), zap.ErrorOutput(w))
	logger := zap.New(core, options...)
	return logger, atomicLogLevel, nil
}

// _getZapFieldsFromGenerics parses and returns zap.Field from generics and ignores others
func _getZapFieldsFromGenerics(a ...any) []zap.Field {
	var res []zap.Field
	for _, val := range a {
		if field, ok := val.(zap.Field); ok {
			res = append(res, field)
		}
	}

	return res
}
