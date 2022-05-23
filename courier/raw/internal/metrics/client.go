package metrics

import (
	"time"

	"lib.kevinlin.info/aperture"
	"lib.kevinlin.info/aperture/lib"
	"lib.kevinlin.info/aperture/protocol"

	"courier/internal/config"
	"courier/internal/meta"
	"courier/internal/util"
)

var (
	// Client is a globally available Aperture statsd client singleton.
	Client aperture.Statsd = lib.NewNoopStatsd()
)

// initClient statefully initializes the global metrics client.
func initClient(cfg *config.Config) error {
	var err error

	serializers := map[string]protocol.LineSerializer{
		"":       protocol.DefaultStatsdLineSerializer,
		"statsd": protocol.DefaultStatsdLineSerializer,
		"influx": protocol.DefaultInfluxLineSerializer,
	}

	Client, err = aperture.NewClient(&aperture.Config{
		Address:                     cfg.Application.Metrics.Address,
		Prefix:                      cfg.Application.Metrics.Prefix,
		Proxy:                       cfg.Application.Metrics.Proxy,
		Serializer:                  serializers[cfg.Application.Metrics.Serializer],
		AsyncQueueSize:              1024,
		TransportProbeInterval:      10 * time.Second,
		LazyTransportInitialization: true,
		DefaultTags: map[string]interface{}{
			"instance": cfg.Application.Instance,
			"version":  meta.Version,
		},
	})
	if err != nil {
		return &util.Error{
			Namespace:    "util",
			Message:      "failed to initialize metrics client",
			StackedError: err,
		}
	}

	return nil
}
