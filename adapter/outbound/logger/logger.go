package logger

import (
	"go.uber.org/fx"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var Module = fx.Module(
	"logger",
	fx.Provide(NewLogger),
	fx.Provide(NewSugarLogger),
)

func NewLogger() (*zap.Logger, error) {
	logger, _ := zap.Config{
		Encoding:          "json",
		Level:             zap.NewAtomicLevelAt(zap.DebugLevel),
		OutputPaths:       []string{"stdout"},
		DisableStacktrace: true,
		EncoderConfig: zapcore.EncoderConfig{
			MessageKey:     "message",
			LevelKey:       "severity",
			TimeKey:        "time",
			NameKey:        "logger",
			CallerKey:      "caller",
			LineEnding:     zapcore.DefaultLineEnding,
			EncodeLevel:    zapcore.LowercaseLevelEncoder,
			EncodeTime:     zapcore.TimeEncoderOfLayout("2006-01-02T15:04:05"),
			EncodeDuration: zapcore.SecondsDurationEncoder,
			EncodeCaller:   zapcore.ShortCallerEncoder,
		},
	}.Build()

	return logger, nil
}

func NewSugarLogger(logger *zap.Logger) *zap.SugaredLogger {
	return logger.Sugar()
}
