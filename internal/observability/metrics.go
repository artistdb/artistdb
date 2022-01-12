package observability

import (
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"

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
	Metrics = newCollector()

	serviceCol := prometheus.NewGaugeFunc(
		prometheus.GaugeOpts{
			Name: "service_info",
			ConstLabels: prometheus.Labels{
				"service": internal.Name,
				"version": internal.Version,
			},
		},
		func() float64 {
			return 1
		},
	)

	cols := []prometheus.Collector{
		collectors.NewBuildInfoCollector(),
		Metrics,
		serviceCol,
	}

	prometheus.MustRegister(cols...)
}

var Metrics *collector

// collector holds all metrics.
type collector struct {
	dbCommandDuration *prometheus.HistogramVec
	dbCommandErrors   *prometheus.CounterVec

	serverRequestDurations *prometheus.HistogramVec
	serverRequestSize      *prometheus.HistogramVec
	serverResponseSize     *prometheus.HistogramVec

	dbObjectsChanged   *prometheus.CounterVec
	dbObjectsRetrieved *prometheus.CounterVec
	dbObjectErrors     *prometheus.CounterVec
}

func (c *collector) Describe(ch chan<- *prometheus.Desc) {
	prometheus.DescribeByCollect(c, ch)
}

func (c *collector) Collect(ch chan<- prometheus.Metric) {
	c.dbCommandDuration.Collect(ch)
	c.dbCommandErrors.Collect(ch)

	c.serverRequestDurations.Collect(ch)
	c.serverRequestSize.Collect(ch)
	c.serverResponseSize.Collect(ch)

	c.dbObjectsChanged.Collect(ch)
	c.dbObjectsRetrieved.Collect(ch)
	c.dbObjectErrors.Collect(ch)
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
		dbCommandErrors: prometheus.NewCounterVec(prometheus.CounterOpts{
			Namespace:   internal.Name,
			Subsystem:   subSystemDB,
			Name:        "command_errors_total",
			Help:        "Total number of errors that occurred from DB commands.",
			ConstLabels: nil,
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
		dbObjectsChanged: prometheus.NewCounterVec(prometheus.CounterOpts{
			Namespace:   internal.Name,
			Subsystem:   subSystemDB,
			Name:        "objects_changed_total",
			Help:        "Total number of objects created, updated or deleted",
			ConstLabels: nil,
		}, []string{"entity", "operation"}),
		dbObjectsRetrieved: prometheus.NewCounterVec(prometheus.CounterOpts{
			Namespace:   internal.Name,
			Subsystem:   subSystemDB,
			Name:        "objects_retrieved_total",
			Help:        "Total number of objects retrieved",
			ConstLabels: nil,
		}, []string{"entity"}),
		dbObjectErrors: prometheus.NewCounterVec(prometheus.CounterOpts{
			Namespace:   internal.Name,
			Subsystem:   subSystemDB,
			Name:        "object_errors_total",
			Help:        "Total number of errors that occurred during interacting with objects",
			ConstLabels: nil,
		}, []string{"entity", "operation"}),
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

func (c *collector) TrackObjectsChanged(amount int, entity string, operation string) {
	c.dbObjectsChanged.WithLabelValues(entity, operation).Add(float64(amount))
}

func (c *collector) TrackObjectsRetrieved(amount int, entity string) {
	c.dbObjectsRetrieved.WithLabelValues(entity).Add(float64(amount))
}

func (c *collector) TrackObjectError(entity string, operation string) {
	c.dbObjectErrors.WithLabelValues(entity, operation).Inc()
}

func (c *collector) TrackCommandError(commandName string) {
	c.dbCommandErrors.WithLabelValues(commandName).Inc()
}
