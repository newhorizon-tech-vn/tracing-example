package prometheus

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

type CounterVec struct {
	counter *prometheus.CounterVec
}

func NewCounterVec(name, description string, labelNames []string) *CounterVec {
	rs := &CounterVec{}

	rs.counter = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: name,
		Help: description,
	},
		labelNames,
	)

	return rs
}

func (c *CounterVec) Increase(lvs ...string) {
	c.counter.WithLabelValues(lvs...).Inc()
}
