package metrics

import (
	"log"
	"time"

	"github.com/prometheus/client_golang/prometheus"

	"github.com/obitech/artist-db/internal"
)

const (
	subSystemDB = "database"
)

func init() {
	Collector = newCollector()

	if err := prometheus.Register(Collector); err != nil {
		log.Fatalf("registering collector failed: %v", err)
	}
}

var Collector *collector

// collector holds all metrics.
type collector struct {
	dbCommandDuration *prometheus.HistogramVec
}

func (c *collector) Describe(ch chan<- *prometheus.Desc) {
	prometheus.DescribeByCollect(c, ch)
}

func (c *collector) Collect(ch chan<- prometheus.Metric) {
	c.dbCommandDuration.Collect(ch)
}

func newCollector() *collector {
	return &collector{
		dbCommandDuration: prometheus.NewHistogramVec(prometheus.HistogramOpts{
			Namespace: internal.Name,
			Subsystem: subSystemDB,
			Name:      "command_duration_seconds",
			Help:      "Observation of command durations against the database.",
			Buckets:   prometheus.ExponentialBuckets(0.05, 2, 10),
		}, []string{"command"}),
	}
}

func (c *collector) ObserveCommandDuration(commandName string, duration time.Duration) {
	c.dbCommandDuration.WithLabelValues(commandName).Observe(duration.Seconds())
}
