package log

import (
	"log"
	"sync"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"courier/internal/config"
	"courier/internal/meta"
	"courier/internal/util"
)

var (
	// StdWarnLogger is a standard library logger that redirects to Zap at zap.WarnLevel.
	StdWarnLogger *log.Logger
	// StdErrorLogger is a standard library logger that redirects to Zap at zap.ErrorLevel.
	StdErrorLogger *log.Logger
)

var (
	mutex sync.Mutex
)

// Init statefully initializes the globally available logging subsystem.
func Init(cfg *config.Config) error {
	mutex.Lock()
	defer mutex.Unlock()

	level := zap.NewAtomicLevel()
	if err := level.UnmarshalText([]byte(cfg.Application.Log.Level)); err != nil {
		return &util.Error{
			Namespace: "log",
			Message:   "failed to parse log level",
			Tags: map[string]interface{}{
				"level": cfg.Application.Log.Level,
			},
			StackedError: err,
		}
	}

	encoder := zap.NewProductionEncoderConfig()
	encoder.EncodeDuration = zapcore.StringDurationEncoder
	encoder.FunctionKey = "func"

	logCfg := zap.Config{
		Level:            level,
		Encoding:         "json",
		EncoderConfig:    encoder,
		OutputPaths:      cfg.Application.Log.Outputs,
		ErrorOutputPaths: cfg.Application.Log.Outputs,
		InitialFields: map[string]interface{}{
			"instance": cfg.Application.Instance,
			"version":  meta.Version,
		},
	}

	logger, err := logCfg.Build()
	if err != nil {
		return &util.Error{
			Namespace:    "log",
			Message:      "failed to create logger from config",
			StackedError: err,
		}
	}

	zap.ReplaceGlobals(logger)

	StdWarnLogger, _ = zap.NewStdLogAt(logger, zap.WarnLevel)
	StdErrorLogger, _ = zap.NewStdLogAt(logger, zap.ErrorLevel)

	return nil
}
