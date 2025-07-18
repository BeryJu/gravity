package extconfig

import (
	"beryju.io/gravity/pkg/extconfig/log_iml"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func (e *ExtConfig) Logger() *zap.Logger {
	return e.logger
}

func (e *ExtConfig) intLog() *zap.Logger {
	return e.Logger().Named("extconfig")
}

func (e *ExtConfig) BuildLogger() *zap.Logger {
	l, err := zapcore.ParseLevel(e.LogLevel)
	if err != nil {
		l = zapcore.InfoLevel
	}
	if e.Debug {
		l = zapcore.DebugLevel
	}
	return e.BuildLoggerWithLevel(l)
}

func (e *ExtConfig) BuildLoggerWithLevel(l zapcore.Level) *zap.Logger {
	config := zap.Config{
		Encoding:         "json",
		Development:      e.Debug,
		OutputPaths:      []string{"stdout"},
		ErrorOutputPaths: []string{"stderr"},
		EncoderConfig:    zap.NewProductionEncoderConfig(),
	}
	config.Level = zap.NewAtomicLevelAt(l)
	config.DisableCaller = !e.Debug
	if e.Debug {
		config.Encoding = "console"
		config.EncoderConfig = zap.NewDevelopmentEncoderConfig()
		config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	}
	if CI() {
		config.EncoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	}
	config.EncoderConfig.EncodeDuration = zapcore.MillisDurationEncoder
	core, err := config.Build()
	if err != nil {
		panic(err)
	}
	hookedCore := &log_iml.ZapCore{
		Core: core.Core(),
	}
	log := zap.New(hookedCore)
	log.Debug("test")
	return log.With(
		zap.String("instance", e.Instance.Identifier),
		zap.String("version", FullVersion()),
	)
}
