package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "dockerhub-rate-limit-exporter",
	Short: "DockerHub Rate Limit Exporter is a tool to fetch and scrape DockerHub pull limits",
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}
