package cmd

import (
	"time"

	"github.com/spf13/cobra"
	"github.com/tcaty/dockerhub-rate-limit-exporter/internal/exporter"
	"github.com/tcaty/dockerhub-rate-limit-exporter/internal/httpserver"
	"github.com/tcaty/dockerhub-rate-limit-exporter/internal/scraper"
	"github.com/tcaty/dockerhub-rate-limit-exporter/pkg/utils"
)

type HttpServer struct {
	Url         string
	MetricsPath string
}

type ScrapeFlags struct {
	ScrapeInterval time.Duration
	HttpServer     HttpServer
	DockerHub      DockerHub
}

var scrapeFlags = &ScrapeFlags{}
var scrapeCmd = &cobra.Command{
	Use:   "scrape",
	Short: "Scrape rate limit data from DockerHub and expose it in Prometheus format.",
	Run: func(cmd *cobra.Command, args []string) {
		initLogger()

		scraper := scraper.New(rootFlags.DockerHub.Repository, rootFlags.DockerHub.Username, scrapeFlags.DockerHub.Password)
		httpServer := httpserver.New(scrapeFlags.HttpServer.Url, scrapeFlags.HttpServer.MetricsPath)
		exporter := exporter.New(scraper, httpServer)

		if err := exporter.Run(scrapeFlags.ScrapeInterval); err != nil {
			utils.HandleFatal(err, "error occured while running rate limit exporter")
		}
	},
}

func init() {
	rootCmd.AddCommand(scrapeCmd)

	scrapeCmd.PersistentFlags().DurationVarP(&scrapeFlags.ScrapeInterval, "scrape-interval", "i", time.Second*15, "interval to scrape DockerHub rate limit")
	scrapeCmd.PersistentFlags().StringVarP(&scrapeFlags.HttpServer.Url, "url", "l", "0.0.0.0:8080", "http server url to run on")
	scrapeCmd.PersistentFlags().StringVarP(&scrapeFlags.HttpServer.MetricsPath, "metrics-path", "m", "/metrics", "path to export metrics")
}
