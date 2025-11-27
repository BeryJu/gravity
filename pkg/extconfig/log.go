package extconfig

import (
	"strings"

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
	return e.BuildLoggerWithLevel(e.LogLevelFor("root"))
}

func (e *ExtConfig) LogLevelFor(role string) zapcore.Level {
	rawLevel := ""
	if e.Debug {
		rawLevel = "debug"
	}
	if role == "root" {
		rawLevel = strings.SplitN(e.LogLevel, ",", 2)[0]
	} else {
		for _, pair := range strings.Split(e.LogLevel, ",") {
			if !strings.Contains(pair, "=") {
				continue
			}
			kv := strings.SplitN(pair, "=", 2)
			if kv[0] == role {
				rawLevel = kv[1]
			}
		}
	}
	l, err := zapcore.ParseLevel(rawLevel)
	if err != nil {
		l = zapcore.InfoLevel
	}
	return l
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
	return log.With(
		zap.String("instance", e.Instance.Identifier),
		zap.String("version", FullVersion()),
	)
}

func SetLevel(lvl zapcore.Level) zap.Option {
	return zap.WrapCore(func(c zapcore.Core) zapcore.Core {
		return &LevelLogger{
			Core:  c,
			Level: lvl,
		}
	})
}

type LevelLogger struct {
	zapcore.Core
	Level zapcore.Level
}

func (ll *LevelLogger) Enabled(l zapcore.Level) bool {
	return ll.Level.Enabled(l)
}

func (ll *LevelLogger) Check(ent zapcore.Entry, ce *zapcore.CheckedEntry) *zapcore.CheckedEntry {
	if ll.Enabled(ent.Level) {
		return ce.AddCore(ent, ll)
	}
	return nil
}

func (ll *LevelLogger) With(fields []zapcore.Field) zapcore.Core {
	return &LevelLogger{
		Core:  ll.Core.With(fields),
		Level: ll.Level,
	}
}
