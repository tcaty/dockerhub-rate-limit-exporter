package cmd

import (
	"github.com/spf13/cobra"
	"github.com/tcaty/dockerhub-rate-limit-exporter/internal/scraper"
	"github.com/tcaty/dockerhub-rate-limit-exporter/pkg/utils"
)

var fetchCmd = &cobra.Command{
	Use:   "fetch",
	Short: "Fetch DockerHub rate limits and print them to console.",
	Run: func(cmd *cobra.Command, args []string) {
		initLogger()

		scraper := scraper.New(
			rootFlags.DockerHub.Repository,
			rootFlags.DockerHub.Username,
			rootFlags.DockerHub.Password,
		)

		if err := scraper.Fetch(); err != nil {
			utils.HandleFatal(err, "error occured while fetching rate limits")
		}
	},
}

func init() {
	rootCmd.AddCommand(fetchCmd)
}
