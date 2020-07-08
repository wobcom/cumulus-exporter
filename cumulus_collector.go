package main

import (
	"sync"

	"gitlab.com/wobcom/cumulus-exporter/collector"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/common/log"
)

type cumulusCollector struct{}

func newCumulusCollector() *cumulusCollector {
	return &cumulusCollector{}
}

func (*cumulusCollector) Describe(ch chan<- *prometheus.Desc) {
	for _, collector := range enabledCollectors {
		collector.Describe(ch)
	}
}

func (*cumulusCollector) Collect(ch chan<- prometheus.Metric) {
	waitGroup := &sync.WaitGroup{}
	waitGroup.Add(len(enabledCollectors))

	for _, collector := range enabledCollectors {
		go runCollector(collector, waitGroup, ch)
	}

	waitGroup.Wait()
}

func runCollector(collector collector.Collector, waitGroup *sync.WaitGroup, ch chan<- prometheus.Metric) {
	defer waitGroup.Done()

	errorChan := make(chan error)
	doneChan := make(chan struct{})

	go collector.Collect(ch, errorChan, doneChan)

	for {
		select {
		case err := <-errorChan:
			log.Errorf("Error running collector %s: %v", collector.Name(), err)
		case <-doneChan:
			return
		}
	}
}
