package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
)

func RegisterCollectors() {

	prometheus.Register(
		collectors.NewGoCollector(),
	)

	prometheus.Register(
		collectors.NewProcessCollector(
			collectors.ProcessCollectorOpts{},
		),
	)
}
