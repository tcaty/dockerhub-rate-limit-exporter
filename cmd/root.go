package cmd

import (
	"time"

	"github.com/spf13/cobra"
)

type Flags struct {
	ScrapeInterval time.Duration
	Repository     string
	MetricsPath    string
	HttpServerUrl  string
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
	rootCmd.PersistentFlags().DurationVarP(&rootFlags.ScrapeInterval, "scrape-interval", "s", time.Second*15, "interval to scrape dockerhub rate limit")
	rootCmd.PersistentFlags().StringVarP(&rootFlags.Repository, "repository", "r", "ratelimitpreview/test", "dockerhub repository to scrape")
	rootCmd.PersistentFlags().StringVarP(&rootFlags.MetricsPath, "metrics-path", "m", "/metrics", "path to export metrics")
	rootCmd.PersistentFlags().StringVarP(&rootFlags.HttpServerUrl, "url", "u", "0.0.0.0:8080", "http server url to run on")
}
