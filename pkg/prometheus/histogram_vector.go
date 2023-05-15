package prometheus

import (
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

type HistogramVector struct {
	histogram *prometheus.HistogramVec
}

func NewHistogramVector(name, description string, labelNames []string) *HistogramVector {
	rs := &HistogramVector{}
	rs.histogram = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Name: name,
		Help: description,
	},
		labelNames,
	)
	return rs
}

func (h *HistogramVector) Observe(start, end time.Time, lvs ...string) {
	h.histogram.WithLabelValues(lvs...).Observe(end.Sub(start).Seconds())
}
