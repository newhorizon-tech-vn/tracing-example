package prometheus

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

type Counter struct {
	counter prometheus.Counter
}

func NewCounter(name, description string) *Counter {
	rs := &Counter{}
	rs.counter = promauto.NewCounter(prometheus.CounterOpts{
		Name: name,
		Help: description,
	})

	return rs
}

func (c *Counter) Increase() {
	c.counter.Inc()
}
