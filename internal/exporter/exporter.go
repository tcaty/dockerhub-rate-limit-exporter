package exporter

import (
	"fmt"
	"net/http"
	"time"
)

type DockerHub interface {
	FetchHeaders() (http.Header, error)
}

type RateLimit interface {
	Init(http.Header) error
	Update(http.Header) error
}

type Exporter struct {
	DockerHub DockerHub
	RateLimit RateLimit
}

func New(dockerhub DockerHub, rateLimit RateLimit) *Exporter {
	return &Exporter{
		DockerHub: dockerhub,
		RateLimit: rateLimit,
	}
}

func (e *Exporter) Run(scrapeInterval time.Duration) error {
	headers, err := e.DockerHub.FetchHeaders()
	if err != nil {
		// TODO: write to chanel
		fmt.Println(err)
	}

	if err := e.RateLimit.Init(headers); err != nil {
		// TODO: write to chanel
		fmt.Println(err)
	}

	go (func() {
		headers, err := e.DockerHub.FetchHeaders()
		if err != nil {
			// TODO: write to chanel
			fmt.Println(err)
		}

		if err := e.RateLimit.Update(headers); err != nil {
			// TODO: write to chanel
			fmt.Println(err)
		}

		time.Sleep(scrapeInterval)
	})()

	return nil
}
