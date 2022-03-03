package metrics

import (
	"sync"

	"courier/internal/config"
)

var (
	mutex sync.Mutex
)

// Init statefully initializes the globally available metrics subsystem.
func Init(cfg *config.Config) error {
	mutex.Lock()
	defer mutex.Unlock()

	if err := initClient(cfg); err != nil {
		return err
	}

	initDefaultHeartbeats(cfg)

	return nil
}
