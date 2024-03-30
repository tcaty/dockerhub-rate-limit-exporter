package exporter

import (
	"fmt"
	"time"
)

type MetaData struct {
	Host     string
	Username string
}

type RateLimitData struct {
	Total     float64
	Remaining float64
}

type DockerHub interface {
	FetchMetaData() (*MetaData, error)
	FetchRateLimitData() (*RateLimitData, error)
}

type RateLimit interface {
	Init(*MetaData)
	Update(*RateLimitData)
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
	metaData, err := e.DockerHub.FetchMetaData()
	if err != nil {
		return fmt.Errorf("could not login to dockerhub: %v", err)
	}

	e.RateLimit.Init(metaData)

	go (func() {
		rateLimitData, err := e.DockerHub.FetchRateLimitData()
		if err != nil {
			// TODO: write to chanel
			fmt.Println(err)
		}

		e.RateLimit.Update(rateLimitData)

		time.Sleep(scrapeInterval)
	})()

	return nil
}
