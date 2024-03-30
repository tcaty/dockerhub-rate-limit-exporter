package ratelimit

import (
	"github.com/tcaty/dockerhub-rate-limit-exporter/internal/exporter"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

type RateLimit struct {
	Total     prometheus.Gauge
	Remaining prometheus.Gauge
}

func New() *RateLimit {
	return &RateLimit{}
}

func (rl *RateLimit) Init(metaData *exporter.MetaData) {
	labels := map[string]string{
		"host":     metaData.Host,
		"username": "anonymus",
	}

	rl.Total = promauto.NewGauge(prometheus.GaugeOpts{
		Name:        "dockerhub_rate_limit_total",
		Help:        "RateLimit-Limit",
		ConstLabels: labels,
	})
	rl.Remaining = promauto.NewGauge(prometheus.GaugeOpts{
		Name:        "dockerhub_rate_limit_remaining",
		Help:        "RateLimit-Remaining",
		ConstLabels: labels,
	})
}

func (rl *RateLimit) Update(rateLimitData *exporter.RateLimitData) {
	rl.Total.Set(rateLimitData.Total)
	rl.Remaining.Set(rateLimitData.Remaining)
}
