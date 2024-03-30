package cmd

import (
	"time"

	"github.com/spf13/cobra"
)

// TODO: refactor
type Flags struct {
	ScrapeInterval    time.Duration
	Repository        string
	MetricsPath       string
	HttpServerUrl     string
	DockerHubUsername string
	DockerHubPassword string
}

var (
	rootFlags = &Flags{}
	rootCmd   = &cobra.Command{}
)

func Execute() (*Flags, error) {
	if err := rootCmd.Execute(); err != nil {
		return nil, err
	}
	return rootFlags, nil
}

func init() {
	rootCmd.PersistentFlags().DurationVarP(&rootFlags.ScrapeInterval, "scrape-interval", "s", time.Second*15, "interval to scrape DockerHub rate limit")
	rootCmd.PersistentFlags().StringVarP(&rootFlags.Repository, "repository", "r", "ratelimitpreview/test", "DockerHub repository to scrape")
	rootCmd.PersistentFlags().StringVarP(&rootFlags.MetricsPath, "metrics-path", "m", "/metrics", "path to export metrics")
	rootCmd.PersistentFlags().StringVarP(&rootFlags.HttpServerUrl, "url", "l", "0.0.0.0:8080", "http server url to run on")
	rootCmd.PersistentFlags().StringVarP(&rootFlags.DockerHubUsername, "username", "u", "", "DockerHub account username")
	rootCmd.PersistentFlags().StringVarP(&rootFlags.DockerHubPassword, "password", "p", "", "DockerHub account password")
}
