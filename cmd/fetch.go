package cmd

import (
	"github.com/spf13/cobra"
	"github.com/tcaty/dockerhub-rate-limit-exporter/internal/scraper"
	"github.com/tcaty/dockerhub-rate-limit-exporter/pkg/utils"
)

type FetchFlags struct {
	DockerHub DockerHub
}

var fetchFlags = &FetchFlags{}
var fetchCmd = &cobra.Command{
	Use:   "fetch",
	Short: "Fetch DockerHub rate limits and print them to console.",
	Run: func(cmd *cobra.Command, args []string) {
		scraper := scraper.New(
			fetchFlags.DockerHub.Repository,
			fetchFlags.DockerHub.Username,
			fetchFlags.DockerHub.Password,
		)
		if err := scraper.Fetch(); err != nil {
			utils.HandleFatal(err, "error occured while fetching rate limits")
		}
	},
}

func init() {
	rootCmd.AddCommand(fetchCmd)

	fetchCmd.PersistentFlags().StringVarP(&fetchFlags.DockerHub.Repository, "repository", "r", "ratelimitpreview/test", "DockerHub repository to fetch")
	fetchCmd.PersistentFlags().StringVarP(&fetchFlags.DockerHub.Username, "username", "u", "", "DockerHub account username (default \"\")")
	fetchCmd.PersistentFlags().StringVarP(&fetchFlags.DockerHub.Password, "password", "p", "", "DockerHub account password (default \"\")")
}
