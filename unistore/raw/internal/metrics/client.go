package metrics

import (
	"fmt"
	"sync"
	"time"

	"lib.kevinlin.info/aperture"
	"lib.kevinlin.info/aperture/lib"
	"lib.kevinlin.info/aperture/protocol"

	"unistore/internal/config"
	"unistore/internal/meta"
)

var (
	// Client is a globally available Aperture statsd client singleton.
	Client aperture.Statsd = lib.NewNoopStatsd()
)

var (
	mutex sync.Mutex
)

// Init statefully initializes the global metrics client.
func Init(cfg *config.Metrics) error {
	var err error

	mutex.Lock()
	defer mutex.Unlock()

	serializers := map[string]protocol.LineSerializer{
		"":       protocol.DefaultStatsdLineSerializer,
		"statsd": protocol.DefaultStatsdLineSerializer,
		"influx": protocol.DefaultInfluxLineSerializer,
	}

	Client, err = aperture.NewClient(&aperture.Config{
		Address:                     cfg.Address,
		Prefix:                      cfg.Prefix,
		Proxy:                       cfg.Proxy,
		Serializer:                  serializers[cfg.Serializer],
		AsyncQueueSize:              1024,
		TransportProbeInterval:      10 * time.Second,
		LazyTransportInitialization: true,
		DefaultTags: map[string]interface{}{
			"version": meta.Version,
		},
	})
	if err != nil {
		return fmt.Errorf("metrics: error creating Aperture client: err=%v", err)
	}

	heartbeats := []lib.Heartbeat{
		lib.NewUptimeHeartbeat(Client),
		lib.NewRuntimeStatsHeartbeat(Client),
		lib.NewResourceUsageHeartbeat(Client),
	}

	for _, hb := range heartbeats {
		lib.RegisterHeartbeat(hb, lib.DefaultHeartbeatInterval, lib.DefaultHeartbeatJitter)
	}

	return nil
}
