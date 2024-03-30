package main

import (
	"github.com/tcaty/dockerhub-rate-limit-exporter/cmd"
	"github.com/tcaty/dockerhub-rate-limit-exporter/internal/dockerhub"
	"github.com/tcaty/dockerhub-rate-limit-exporter/internal/exporter"
	"github.com/tcaty/dockerhub-rate-limit-exporter/internal/httpserver"
	"github.com/tcaty/dockerhub-rate-limit-exporter/internal/ratelimit"
	"github.com/tcaty/dockerhub-rate-limit-exporter/pkg/utils"
)

func main() {
	flags, err := cmd.Execute()
	if err != nil {
		utils.HandleFatal(err, "error occured during command line flags initialization")
	}

	dockerhub := dockerhub.New(flags.Repository, flags.DockerHubUsername, flags.DockerHubPassword)
	ratelimit := ratelimit.New()
	exporter := exporter.New(dockerhub, ratelimit)

	if err := exporter.Run(flags.ScrapeInterval); err != nil {
		utils.HandleFatal(err, "error occured while running exporter")
	}

	httpServer := httpserver.New(flags.HttpServerUrl, flags.MetricsPath)
	if err := httpServer.Run(); err != nil {
		utils.HandleFatal(err, "error occured while running http server")
	}
}
