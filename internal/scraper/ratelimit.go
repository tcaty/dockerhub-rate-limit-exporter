package scraper

import (
	"fmt"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

type RateLimit struct {
	Total     prometheus.Gauge
	Remaining prometheus.Gauge
}

func NewRateLimit(metaData *MetaData) *RateLimit {
	labels := map[string]string{
		"host":     metaData.Host,
		"username": metaData.Username,
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

	return &RateLimit{
		Total:     total,
		Remaining: remaining,
	}
}

func (rl *RateLimit) Update(rateLimitData *RateLimitData) {
	rl.Total.Set(rateLimitData.Total)
	rl.Remaining.Set(rateLimitData.Remaining)
	fmt.Println("updated")
}
