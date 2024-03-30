package main

import (
	"github.com/tcaty/dockerhub-rate-limit-exporter/cmd"
	"github.com/tcaty/dockerhub-rate-limit-exporter/internal/exporter"
	"github.com/tcaty/dockerhub-rate-limit-exporter/internal/httpserver"
	"github.com/tcaty/dockerhub-rate-limit-exporter/internal/scraper"
	"github.com/tcaty/dockerhub-rate-limit-exporter/pkg/utils"
)

func main() {
	flags, err := cmd.Execute()

	if err != nil {
		utils.HandleFatal(err, "error occured during command line flags initialization")
	}

	scraper := scraper.New(*flags)
	httpServer := httpserver.New(*flags)
	exporter := exporter.New(scraper, httpServer)
	// TODO: split to subcmds
	scraper.Fetch()
	if err := exporter.Run(flags.ScrapeInterval); err != nil {
		utils.HandleFatal(err, "error occured during running exporter")
	}
}
