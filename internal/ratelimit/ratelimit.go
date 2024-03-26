package ratelimit

import (
	"net/http"

	"github.com/tcaty/dockerhub-rate-limit-exporter/pkg/utils"

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

func (rl *RateLimit) Init(headers http.Header) error {
	host, err := utils.ParseHeader(headers, "Docker-Ratelimit-Source")
	if err != nil {
		return err
	}

	labels := map[string]string{
		"host":     host,
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

	return nil
}

func (rl *RateLimit) Update(headers http.Header) error {
	limit, err := parseRateLimitHeader(headers, "Limit")
	if err != nil {
		return err
	}

	remaining, err := parseRateLimitHeader(headers, "Remaining")
	if err != nil {
		return err
	}

	rl.Total.Set(limit)
	rl.Remaining.Set(remaining)

	return nil
}
