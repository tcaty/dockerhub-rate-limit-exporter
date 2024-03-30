package exporter

import (
	"fmt"
	"time"

	"github.com/tcaty/dockerhub-rate-limit-exporter/internal/httpserver"
	"github.com/tcaty/dockerhub-rate-limit-exporter/internal/scraper"
)

type Exporter struct {
	scraper    *scraper.Scraper
	httpServer *httpserver.HttpServer
}

func New(scraper *scraper.Scraper, httpServer *httpserver.HttpServer) *Exporter {
	return &Exporter{
		scraper:    scraper,
		httpServer: httpServer,
	}
}

func (e *Exporter) Run(scrapeInterval time.Duration) error {
	httpServerErrCh := make(chan error)
	scraperErrCh := make(chan error)

	go func() {
		if err := e.httpServer.Run(); err != nil {
			httpServerErrCh <- err
		}
	}()

	go func() {
		if err := e.scraper.Scrape(scrapeInterval); err != nil {
			scraperErrCh <- err
		}
	}()

	select {
	case err := <-httpServerErrCh:
		return fmt.Errorf("error occured while running httpServer: %v", err)
	case err := <-scraperErrCh:
		return fmt.Errorf("error occured while running scraper: %v", err)
	}
}
