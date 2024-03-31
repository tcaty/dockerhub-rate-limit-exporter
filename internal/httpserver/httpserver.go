package httpserver

import (
	"log/slog"
	"net/http"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type HttpServer struct {
	url         string
	metricsPath string
}

func New(url string, metricsPath string) *HttpServer {
	return &HttpServer{
		url:         url,
		metricsPath: metricsPath,
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
