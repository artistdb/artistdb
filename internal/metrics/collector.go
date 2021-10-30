package metrics

import "github.com/prometheus/client_golang/prometheus"

// Collector holds all metrics.
type Collector struct{}

func (c *Collector) Describe(descs chan<- *prometheus.Desc) {
	prometheus.DescribeByCollect(c, descs)
}

func (c *Collector) Collect(metrics chan<- prometheus.Metric) {}

func NewCollector() *Collector {
	return &Collector{}
}
