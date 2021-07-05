// This example program demonstrates basic usage of the Aperture statsd client with minimal
// additional configuration customization.

package main

import (
	"time"

	"lib.kevinlin.info/aperture"
)

func main() {
	// Report some sample metrics to a UDP statsd server listening on localhost:8125.

	client, err := aperture.NewClient(&aperture.Config{
		Address: "udp://localhost:8125",
		Prefix:  "aperture",
	})
	if err != nil {
		panic(err)
	}

	client.Incr("counter", nil)
	client.Timing("timer", 150*time.Millisecond, map[string]interface{}{"error": false})
	client.Gauge("gauge", 10, nil)

	client.Close()
}
