package main

import (
	"fmt"

	"github.com/tcaty/dockerhub-rate-limit-exporter/cmd"
	"github.com/tcaty/dockerhub-rate-limit-exporter/internal/dockerhub"
	"github.com/tcaty/dockerhub-rate-limit-exporter/internal/exporter"
	"github.com/tcaty/dockerhub-rate-limit-exporter/internal/httpserver"
	"github.com/tcaty/dockerhub-rate-limit-exporter/internal/ratelimit"
)

func main() {
	flags := cmd.Execute()
	fmt.Println(flags)

	dockerhub := dockerhub.New(flags.Repository)
	ratelimit := ratelimit.New()
	exporter := exporter.New(dockerhub, ratelimit)
	exporter.Run(flags.ScrapeInterval)

	httpServer := httpserver.New(flags.HttpServerPort, flags.MetricsPath)
	httpServer.Run()
}
