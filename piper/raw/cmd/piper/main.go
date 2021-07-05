package main

import (
	"flag"
	"os"
	"time"

	"github.com/getsentry/sentry-go"
	"go.uber.org/zap"
	"lib.kevinlin.info/aperture"
	"lib.kevinlin.info/aperture/lib"

	"piper/internal/config"
	"piper/internal/meta"
	"piper/internal/relay"
	"piper/internal/tail"
	"piper/internal/util"
)

var (
	cfg     *config.Config
	metrics aperture.Statsd = lib.NewNoopStatsd()
	logger  *zap.Logger
)

var (
	flagConfig    = flag.String("config", os.Getenv("PIPER_CONFIG"), "path to the configuration file on disk")
	flagVerbosity = flag.String("verbosity", zap.InfoLevel.String(), "log level")
)

func init() {
	flag.Parse()

	level := zap.NewAtomicLevel()
	level.UnmarshalText([]byte(*flagVerbosity))
	logCfg := zap.Config{
		Level:            level,
		Encoding:         "json",
		EncoderConfig:    zap.NewProductionEncoderConfig(),
		OutputPaths:      []string{"stderr"},
		ErrorOutputPaths: []string{"stderr"},
	}
	logger, _ = logCfg.Build()

	if *flagConfig == "" {
		logger.Fatal("path to config file must be supplied")
	}
}

func main() {
	var err error

	defer logger.Sync()

	logger.Info(
		"starting piper daemon",
		zap.String("version", meta.Version),
		zap.String("config", *flagConfig),
	)

	cfg, err = config.New(*flagConfig)
	if err != nil {
		logger.Fatal("failed to create config", zap.Error(err))
	}

	if cfg.Piper.MetricsAddress != "" {
		logger.Info(
			"enabling metrics reporting",
			zap.String("prefix", cfg.Piper.MetricsPrefix),
			zap.String("address", cfg.Piper.MetricsAddress),
			zap.String("proxy", cfg.Piper.MetricsProxy),
		)

		tags, err := cfg.Piper.DefaultMetricsTags()
		if err != nil {
			logger.Fatal("failed to enumerate default metrics tags", zap.Error(err))
		}

		metrics, err = aperture.NewClient(&aperture.Config{
			Address:                     cfg.Piper.MetricsAddress,
			Proxy:                       cfg.Piper.MetricsProxy,
			Prefix:                      cfg.Piper.MetricsPrefix,
			TransportProbeInterval:      10 * time.Second,
			TransportKeepaliveInterval:  1 * time.Minute,
			TransportConnectTimeout:     5 * time.Second,
			TransportSendTimeout:        5 * time.Second,
			LazyTransportInitialization: true,
			DefaultTags:                 tags,
		})
		if err != nil {
			logger.Fatal("failed to create metrics client", zap.Error(err))
		}

		metrics = lib.NewAsyncStatsd(metrics)
	}

	if cfg.Piper.SentryDSN != "" {
		logger.Info(
			"enabling error reporting",
			zap.String("sentry_dsn", cfg.Piper.SentryDSN),
		)

		err := sentry.Init(sentry.ClientOptions{
			Dsn:              cfg.Piper.SentryDSN,
			Release:          meta.Version,
			Transport:        sentry.NewHTTPSyncTransport(),
			AttachStacktrace: true,
		})
		if err != nil {
			logger.Fatal("failed to create Sentry client", zap.Error(err))
		}

		defer sentry.Recover()
	}

	for _, cfgRelay := range cfg.Relays {
		logger.Info(
			"registering log stream relay",
			zap.String("name", cfgRelay.Name),
			zap.String("path", cfgRelay.LogFile.Pattern),
			zap.String("address", cfgRelay.KafkaAddress),
			zap.String("topic", cfgRelay.KafkaTopic),
		)

		go func(cfgRelay *config.Relay) {
			for {
				logger.Info(
					"initializing log stream relay",
					zap.String("name", cfgRelay.Name),
					zap.String("path", cfgRelay.LogFile.Pattern),
					zap.String("address", cfgRelay.KafkaAddress),
					zap.String("topic", cfgRelay.KafkaTopic),
				)

				tags := map[string]interface{}{
					"name":  cfgRelay.Name,
					"path":  cfgRelay.LogFile,
					"topic": cfgRelay.KafkaTopic,
				}

				metrics.Incr("relay.open", tags)

				if err := streamLog(cfgRelay); err != nil {
					metrics.Incr("relay.error", tags)
					logger.Error(
						"log stream relay encountered error",
						zap.Error(err),
					)

					hub := sentry.CurrentHub().Clone()
					hub.ConfigureScope(func(scope *sentry.Scope) {
						scope.SetTag("name", cfgRelay.Name)
						scope.SetTag("path", cfgRelay.LogFile.Pattern)
						scope.SetTag("topic", cfgRelay.KafkaTopic)
					})
					hub.CaptureException(err)
				} else {
					logger.Info("log stream relay concluded without error")
				}

				metrics.Incr("relay.close", tags)

				logger.Info(
					"scheduling log stream relay re-initialization",
					zap.String("name", cfgRelay.Name),
					zap.Duration("delay", cfgRelay.ResetDelay.Duration),
				)
				time.Sleep(cfgRelay.ResetDelay.Duration)
			}
		}(cfgRelay)
	}

	<-make(chan interface{})

}

// streamLog opens a tailer on a log file and starts a Kafka producer for new lines.
// The function returns when the tailer completes or encounters an error.
func streamLog(cfgRelay *config.Relay) error {
	// Tailer initialization

	tailOpts := tail.Options{
		BufferLength: cfgRelay.BufferLength,
		SeekPosition: cfgRelay.SeekPosition,
	}
	if len(cfgRelay.Delimiter) == 1 {
		tailOpts.Delimiter = cfgRelay.Delimiter[0]
	}
	if cfgRelay.SeekPosition == "" {
		tailOpts.SeekPosition = config.SeekPositionEnd
	}
	if len(cfgRelay.TagIdentifiers) > 0 {
		tags, err := cfgRelay.Tags()
		if err != nil {
			return err
		}
		tailOpts.Tags = tags
	}

	match, err := cfgRelay.LogFile.Match()
	if err != nil {
		return err
	}

	logFile, err := tail.New(match, tailOpts)
	if err != nil {
		return err
	}
	defer logFile.Close()

	logger.Debug(
		"opened log file for tailing",
		zap.String("path", match),
		zap.Any("options", tailOpts),
	)

	// Kafka producer initialization

	producer, err := relay.NewProducer(
		relay.MessageSerializerRegistry[cfgRelay.Serializer],
		cfgRelay.ProxyAddress,
		cfgRelay.KafkaAddress,
		cfgRelay.KafkaAcks,
		cfgRelay.KafkaRetries,
		cfgRelay.KafkaTimeout.Duration,
	)
	if err != nil {
		return err
	}
	defer producer.Close()

	logger.Debug(
		"established Kafka producer connection",
		zap.String("proxy", cfgRelay.ProxyAddress),
		zap.String("address", cfgRelay.KafkaAddress),
		zap.Int("acks", cfgRelay.KafkaAcks),
		zap.Int("retries", cfgRelay.KafkaRetries),
		zap.Duration("timeout", cfgRelay.KafkaTimeout.Duration),
	)

	// Consume tailed messaged from the log and ship them to Kafka

	messages, err := logFile.Consume()
	if err != nil {
		return err
	}

	logger.Debug("listening for tail events in log stream")

	for message := range messages {
		tags := map[string]interface{}{
			"name":  cfgRelay.Name,
			"path":  cfgRelay.LogFile,
			"topic": cfgRelay.KafkaTopic,
		}

		metrics.Incr("relay.log.tail", tags)
		metrics.Gauge("relay.log.sequence_id", float64(message.SequenceID), tags)

		if cfgRelay.Filter.Regexp != nil && !cfgRelay.Filter.MatchString(message.Line) {
			metrics.Incr("relay.log.skip", tags)
			logger.Debug(
				"skipping message due to absence of filter pattern match",
				zap.Stringer("filter", cfgRelay.Filter),
				zap.String("topic", cfgRelay.KafkaTopic),
				zap.Time("timestamp", message.Timestamp),
				zap.String("path", message.Path),
				zap.String("line", message.Line),
			)
			continue
		}

		for {
			logger.Debug(
				"publishing tail event message",
				zap.String("topic", cfgRelay.KafkaTopic),
				zap.Time("timestamp", message.Timestamp),
				zap.String("path", message.Path),
				zap.String("line", message.Line),
			)

			stopwatch := lib.NewStopwatch()
			publishTags := util.MergeMaps(tags, map[string]interface{}{"success": true})

			size, err := producer.WriteMessage(cfgRelay.KafkaTopic, message)
			if err != nil {
				publishTags["success"] = false
				logger.Error("error publishing tail event message", zap.Error(err))
			}

			metrics.Incr("relay.log.publish.message", publishTags)
			metrics.Timing("relay.log.publish.latency", stopwatch.Elapsed(), publishTags)
			metrics.Timing("relay.log.publish.lag", time.Now().Sub(message.Timestamp), publishTags)
			metrics.Size("relay.log.publish.size", int64(size), publishTags)

			if err == nil {
				break
			}
		}
	}

	return logFile.Error()
}
