package metrics

import (
	"log"
	"time"

	"github.com/prometheus/client_golang/prometheus"

	"github.com/obitech/artist-db/internal"
)

const (
	subSystemDB     = "database"
	subSystemServer = "server"
)

var (
	serverLabels      = []string{"method", "route", "code"}
	serverSizeBuckets = []float64{50, 150, 300, 800, 1_200, 5_000, 8_000, 10_000, 20_000}
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
	dbCommandDuration      *prometheus.HistogramVec
	serverRequestDurations *prometheus.HistogramVec
	serverRequestSize      *prometheus.HistogramVec
	serverResponseSize     *prometheus.HistogramVec
}

func (c *collector) Describe(ch chan<- *prometheus.Desc) {
	prometheus.DescribeByCollect(c, ch)
}

func (c *collector) Collect(ch chan<- prometheus.Metric) {
	c.dbCommandDuration.Collect(ch)

	c.serverRequestDurations.Collect(ch)
	c.serverRequestSize.Collect(ch)
	c.serverResponseSize.Collect(ch)
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
		serverRequestDurations: prometheus.NewHistogramVec(prometheus.HistogramOpts{
			Namespace: internal.Name,
			Subsystem: subSystemServer,
			Name:      "request_duration_seconds",
			Help:      "Latency of HTTP requests.",
		}, serverLabels),
		serverRequestSize: prometheus.NewHistogramVec(prometheus.HistogramOpts{
			Namespace: internal.Name,
			Subsystem: subSystemServer,
			Name:      "request_size_bytes",
			Help:      "Size of HTTP requests.",
			Buckets:   serverSizeBuckets,
		}, serverLabels),
		serverResponseSize: prometheus.NewHistogramVec(prometheus.HistogramOpts{
			Namespace: internal.Name,
			Subsystem: subSystemServer,
			Name:      "response_size_bytes",
			Help:      "Size of HTTP responses.",
			Buckets:   serverSizeBuckets,
		}, serverLabels),
	}
}

func (c *collector) ObserveCommandDuration(commandName string, duration time.Duration) {
	c.dbCommandDuration.WithLabelValues(commandName).Observe(duration.Seconds())
}

func (c *collector) ObserveRequestDuration(method, route, code string, duration time.Duration) {
	c.serverRequestDurations.WithLabelValues(method, route, code).Observe(duration.Seconds())
}

func (c *collector) ObserveRequestSize(method, route, code string, size float64) {
	c.serverRequestDurations.WithLabelValues(method, route, code).Observe(size)
}

func (c *collector) ObserveResponseSize(method, route, code string, size float64) {
	c.serverResponseSize.WithLabelValues(method, route, code).Observe(size)
}
