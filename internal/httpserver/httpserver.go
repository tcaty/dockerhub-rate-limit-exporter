package httpserver

import (
	"log/slog"
	"net/http"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/tcaty/dockerhub-rate-limit-exporter/cmd"
)

type HttpServer struct {
	url         string
	metricsPath string
}

func New(flags cmd.Flags) *HttpServer {
	return &HttpServer{
		url:         flags.HttpServerUrl,
		metricsPath: flags.MetricsPath,
	}
}

func (hs *HttpServer) Run() error {
	http.Handle(hs.metricsPath, promhttp.Handler())

	slog.Info(
		"Starting http server",
		"url", hs.url,
		"path", hs.metricsPath,
	)

	return http.ListenAndServe(hs.url, nil)
}
