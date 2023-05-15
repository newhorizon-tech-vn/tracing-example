package prometheus

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

type Histogram struct {
	histogram prometheus.Histogram
}

func NewHistogram(name, description string) *Histogram {
	rs := &Histogram{}
	rs.histogram = promauto.NewHistogram(
		prometheus.HistogramOpts{
			Name:    name,
			Help:    description,
			Buckets: prometheus.DefBuckets,
		},
	)
	return rs
}

func (h *Histogram) Observe(t float64) {
	h.histogram.Observe(t)
}
