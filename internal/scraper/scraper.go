package scraper

import (
	"fmt"
	"log/slog"
	"time"
)

type metaData struct {
	host     string
	username string
}

type rateLimitData struct {
	total     float64
	remaining float64
}

type Scraper struct {
	dockerHub *dockerHub
}

func New(repository string, username string, password string) *Scraper {
	return &Scraper{
		dockerHub: &dockerHub{
			repository: repository,
			username:   username,
			password:   password,
		},
	}
}

func (s *Scraper) Scrape(interval time.Duration) error {
	metaData, err := s.dockerHub.fetchMetaData()

	if err != nil {
		return fmt.Errorf("could not login to dockerhub: %v", err)
	}

	rateLimit := NewRateLimit(metaData)

	for {
		rateLimitData, err := s.dockerHub.fetchRateLimitData()

		if err != nil {
			return err
		}

		rateLimit.update(rateLimitData)

		slog.Info(
			"rate limit scrape succeeded",
			"host", metaData.host,
			"username", metaData.username,
		)

		time.Sleep(interval)
	}
}

func (s *Scraper) Fetch() error {
	metaData, err := s.dockerHub.fetchMetaData()

	if err != nil {
		return fmt.Errorf("could not login to dockerhub: %v", err)
	}

	rateLimitData, err := s.dockerHub.fetchRateLimitData()

	if err != nil {
		return err
	}

	if s.dockerHub.isAuthenticatedMode() {
		fmt.Println("Mode:", "Authenticated")
		fmt.Println("Username:", metaData.username)
	} else {
		fmt.Println("Mode:", "Anonymous")
	}
	fmt.Println("Host:", metaData.host)
	fmt.Println("RateLimit [Total]:", rateLimitData.total)
	fmt.Println("RateLimit [Remaining]:", rateLimitData.remaining)

	return nil
}
