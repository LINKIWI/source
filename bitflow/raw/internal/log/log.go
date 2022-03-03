package log

import (
	"fmt"
	"io"
	stdlog "log"
	"os"

	"bitflow/internal/config"
)

var (
	defaultFlags = stdlog.Lshortfile | stdlog.Lmsgprefix | stdlog.Ldate | stdlog.Ltime
)

// Context is a container for debug, info, warn, and error Logger instances.
type Context struct {
	Debug *stdlog.Logger
	Info  *stdlog.Logger
	Warn  *stdlog.Logger
	Error *stdlog.Logger
}

// Init statefully initializes the global standard loggers.
func Init(cfg *config.Log) (*Context, error) {
	ctx := &Context{
		Debug: stdlog.New(io.Discard, "debug: ", defaultFlags),
		Info:  stdlog.New(io.Discard, "info: ", defaultFlags),
		Warn:  stdlog.New(io.Discard, "warn: ", defaultFlags),
		Error: stdlog.New(io.Discard, "error: ", defaultFlags),
	}

	if cfg == nil {
		return ctx, nil
	}

	switch cfg.Level {
	case "debug":
		ctx.Debug.SetOutput(os.Stdout)
		fallthrough
	case "info":
		ctx.Info.SetOutput(os.Stdout)
		fallthrough
	case "warn":
		ctx.Warn.SetOutput(os.Stdout)
		fallthrough
	case "error":
		ctx.Error.SetOutput(os.Stdout)
	default:
		return nil, fmt.Errorf("log: unrecognized log level: level=%s", cfg.Level)
	}

	return ctx, nil
}
