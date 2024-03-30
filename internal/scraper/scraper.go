package scraper

import (
	"fmt"
	"time"

	"github.com/tcaty/dockerhub-rate-limit-exporter/cmd"
)

// TODO: make private almost all

type MetaData struct {
	Host     string
	Username string
}

type RateLimitData struct {
	Total     float64
	Remaining float64
}

type Scraper struct {
	DockerHub *DockerHub
}

func New(flags cmd.Flags) *Scraper {
	dockerhub := NewDockerHub(flags)

	return &Scraper{
		DockerHub: dockerhub,
	}
}

func (s *Scraper) Scrape(interval time.Duration) error {
	metaData, err := s.DockerHub.FetchMetaData()

	if err != nil {
		return fmt.Errorf("could not login to dockerhub: %v", err)
	}

	rateLimit := NewRateLimit(metaData)

	for {
		rateLimitData, err := s.DockerHub.FetchRateLimitData()

		if err != nil {
			return err
		}

		rateLimit.Update(rateLimitData)

		time.Sleep(interval)
	}
}

func (s *Scraper) Fetch() error {
	metaData, err := s.DockerHub.FetchMetaData()

	if err != nil {
		return fmt.Errorf("could not login to dockerhub: %v", err)
	}

	rateLimitData, err := s.DockerHub.FetchRateLimitData()

	if err != nil {
		return err
	}

	if s.DockerHub.IsAuthenticatedMode() {
		fmt.Println("Mode:", "Authenticated")
		fmt.Println("Username:", metaData.Username)
	} else {
		fmt.Println("Mode:", "Anonymous")
	}
	fmt.Println("Host:", metaData.Host)
	fmt.Println("RateLimit [Total]:", rateLimitData.Total)
	fmt.Println("RateLimit [Remaining]:", rateLimitData.Remaining)

	return nil
}
