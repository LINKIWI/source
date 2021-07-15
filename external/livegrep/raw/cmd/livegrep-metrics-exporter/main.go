package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"

	"lib.kevinlin.info/aperture"
	"lib.kevinlin.info/aperture/lib"
)

const (
	envStatsdAddr   = "LIVEGREP_METRICS_STATSD_ADDRESS"
	envStatsdPrefix = "LIVEGREP_METRICS_STATSD_PREFIX"
)

var (
	flagStatsdAddr   = flag.String("statsd-address", os.Getenv(envStatsdAddr), "address URI of statsd listener for metrics export")
	flagStatsdPrefix = flag.String("statsd-prefix", os.Getenv(envStatsdPrefix), "optional prefix to apply to all metrics")
	flagStatsdTags   = newStringMapFlag()

	metricPattern = regexp.MustCompile(strings.Join([]string{
		"([\\w\\.]+)", // Metric name (alphabetic characters and dots)
		"\\s",         // Separator
		"(\\d+)",      // Metric value (gauge)
	}, ""))
	metricsDumpPattern = regexp.MustCompile("(?s)== begin metrics ==\\s*(.+)\\s== end metrics ==")
)

func init() {
	flag.Var(flagStatsdTags, "statsd-tag", "statsd tags to include on all emitted metrics")

	flag.Parse()

	if *flagStatsdAddr == "" {
		panic(fmt.Errorf("no statsd target address specified"))
	} else {
		log.Printf("using statsd server: address=%s", *flagStatsdAddr)
	}

	if *flagStatsdPrefix != "" {
		log.Printf("using prefix for all metrics: prefix=%s", *flagStatsdPrefix)
	}
}

func main() {
	log.Println("starting livegrep statsd metrics exporter")

	// Stopwatch to track the end-to-end duration required to export metrics
	stopwatch := lib.NewStopwatch()

	// Create a statsd client
	statsd, err := aperture.NewClient(&aperture.Config{
		Address: *flagStatsdAddr,
		Prefix:  *flagStatsdPrefix,
	})
	if err != nil {
		panic(err)
	}

	// Read full index builder logs from standard input
	indexLogs, err := io.ReadAll(os.Stdin)
	if err != nil {
		panic(fmt.Errorf("failed to read index builder logs from stdin: err=%v", err))
	}

	// Regex-match the metrics dump block
	dump := metricsDumpPattern.FindStringSubmatch(string(indexLogs))
	if len(dump) < 2 {
		panic(fmt.Errorf("failed to parse metrics dump from indexer output"))
	}

	// Regex-match the metric name and value from each line
	metrics := make(map[string]int)
	for _, metricLine := range strings.Split(dump[1], "\n") {
		metric := metricPattern.FindStringSubmatch(metricLine)
		if len(metric) < 3 {
			panic(fmt.Errorf("failed to parse metric name and value: line=%s", metricLine))
		}

		value, err := strconv.Atoi(metric[2])
		if err != nil {
			panic(fmt.Errorf("failed to parse metric value: name=%s value=%s", metric[1], metric[2]))
		}

		metrics[metric[1]] = value
	}

	// Report all parsed gauge metrics to statsd
	for metric, value := range metrics {
		log.Printf("reporting gauge metric: metric=%s value=%d", metric, value)
		statsd.Gauge(metric, float64(value), flagStatsdTags.Values())
	}

	// Report metrics export duration to statsd
	stopwatch.Pause()
	log.Printf("completed metrics export: duration=%v", stopwatch.Elapsed())
	statsd.Timing("export.duration", stopwatch.Elapsed(), flagStatsdTags.Values())
}
