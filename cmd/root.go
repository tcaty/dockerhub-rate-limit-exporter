package cmd

import (
	"os"
	"time"

	"github.com/spf13/cobra"
)

type Flags struct {
	ScrapeInterval time.Duration
	Repository     string
	MetricsPath    string
	HttpServerPort int
}

var (
	rootFlags = &Flags{}
	rootCmd   = &cobra.Command{}
)

func Execute() *Flags {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
	return rootFlags
}

func init() {
	rootCmd.PersistentFlags().DurationVarP(&rootFlags.ScrapeInterval, "scrape-interval", "s", time.Second*15, "interval to scrape dockerhub rate limit")
	rootCmd.PersistentFlags().StringVarP(&rootFlags.Repository, "repository", "r", "ratelimitpreview/test", "dockerhub repository to scrape")
	rootCmd.PersistentFlags().StringVarP(&rootFlags.MetricsPath, "metrics-path", "m", "/metrics", "path to export metrics")
	rootCmd.PersistentFlags().IntVarP(&rootFlags.HttpServerPort, "port", "p", 8080, "http server port to run on")
}
