package scraper

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

type rateLimit struct {
	total     prometheus.Gauge
	remaining prometheus.Gauge
}

func NewRateLimit(metaData *metaData) *rateLimit {
	labels := map[string]string{
		"host":     metaData.host,
		"username": metaData.username,
		"mode":     metaData.mode,
	}

	total := promauto.NewGauge(prometheus.GaugeOpts{
		Name:        "dockerhub_rate_limit_total",
		Help:        "RateLimit-Limit",
		ConstLabels: labels,
	})

	remaining := promauto.NewGauge(prometheus.GaugeOpts{
		Name:        "dockerhub_rate_limit_remaining",
		Help:        "RateLimit-Remaining",
		ConstLabels: labels,
	})

	return &rateLimit{
		total:     total,
		remaining: remaining,
	}
}

func (rl *rateLimit) update(rateLimitData *rateLimitData) {
	rl.total.Set(rateLimitData.total)
	rl.remaining.Set(rateLimitData.remaining)
}
