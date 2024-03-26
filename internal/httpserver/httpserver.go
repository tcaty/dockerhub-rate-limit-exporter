package httpserver

import (
	"fmt"
	"net/http"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type HttpServer struct {
	port        int
	metricsPath string
}

func New(port int, metricsPath string) *HttpServer {
	return &HttpServer{
		port:        port,
		metricsPath: metricsPath,
	}
}

func (s *HttpServer) Run() {
	http.Handle(s.metricsPath, promhttp.Handler())
	http.ListenAndServe(fmt.Sprintf(":%d", s.port), nil)
}
