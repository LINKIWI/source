package log

import (
	"fmt"
	"sync"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"unistore/internal/config"
	"unistore/internal/meta"
)

var (
	mutex sync.Mutex
)

// Init statefully initializes the globally available logging subsystem.
func Init(cfg *config.Log) error {
	mutex.Lock()
	defer mutex.Unlock()

	level := zap.NewAtomicLevel()
	if err := level.UnmarshalText([]byte(cfg.Level)); err != nil {
		return fmt.Errorf("log: error unmarshaling log level: err=%v", err)
	}

	encoder := zap.NewProductionEncoderConfig()
	encoder.EncodeDuration = zapcore.StringDurationEncoder
	encoder.FunctionKey = "func"

	logCfg := zap.Config{
		Level:            level,
		Encoding:         "json",
		EncoderConfig:    encoder,
		OutputPaths:      cfg.Outputs,
		ErrorOutputPaths: cfg.Outputs,
		InitialFields:    map[string]interface{}{"version": meta.Version},
	}

	logger, err := logCfg.Build()
	if err != nil {
		return fmt.Errorf("log: error building logging configuration: err=%v", err)
	}

	zap.ReplaceGlobals(logger)

	return nil
}
