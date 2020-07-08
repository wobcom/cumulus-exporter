package collector

import (
	"github.com/prometheus/client_golang/prometheus"
)

// Collector is an interface, that is implemented by mstpd, asic collector, etc.
type Collector interface {
	Name() string
	Describe(ch chan<- *prometheus.Desc)
	Collect(metrics chan<- prometheus.Metric, errorChan chan error, done chan struct{})
}
