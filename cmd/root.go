package cmd

import (
	"log/slog"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"github.com/tcaty/dockerhub-rate-limit-exporter/pkg/utils"
)

type DockerHub struct {
	Repository string
	Username   string
	Password   string
}

type RootFlags struct {
	LogLevel  string
	DockerHub DockerHub
}

var rootFlags = &RootFlags{}
var rootCmd = &cobra.Command{
	Use:   "dockerhub-rate-limit-exporter",
	Short: "DockerHub Rate Limit Exporter is a tool to fetch and scrape DockerHub pull limits",
}

func Execute() {
	err := rootCmd.Execute()

	if err != nil {
		utils.HandleFatal(err, "error occured while running rootCmd")
	}
}

func init() {
	rootCmd.PersistentFlags().StringVarP(&rootFlags.LogLevel, "log-level", "", slog.LevelInfo.String(), "Slog log level")
	rootCmd.PersistentFlags().StringVarP(&rootFlags.DockerHub.Repository, "repository", "r", "ratelimitpreview/test", "DockerHub repository to fetch")
	rootCmd.PersistentFlags().StringVarP(&rootFlags.DockerHub.Username, "username", "u", "", "DockerHub account username (default \"\")")
	rootCmd.PersistentFlags().StringVarP(&rootFlags.DockerHub.Password, "password", "p", "", "DockerHub account password (default \"\")")

	cobra.OnInitialize(func() {
		viper.AutomaticEnv()
		postInitCommands(rootCmd.Commands())
	})
}

func postInitCommands(commands []*cobra.Command) {
	for _, cmd := range commands {
		presetRequiredFlags(cmd)
		if cmd.HasSubCommands() {
			postInitCommands(cmd.Commands())
		}
	}
}

func presetRequiredFlags(cmd *cobra.Command) {
	viper.BindPFlags(cmd.Flags())
	cmd.Flags().VisitAll(func(f *pflag.Flag) {
		if viper.IsSet(f.Name) && viper.GetString(f.Name) != "" {
			cmd.Flags().Set(f.Name, viper.GetString(f.Name))
		}
	})
}

func initLogger() {
	level := new(slog.LevelVar)

	if err := level.UnmarshalText([]byte(rootFlags.LogLevel)); err != nil {
		utils.HandleFatal(err, "unable to parse log level")
	}

	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: level,
	}))

	slog.SetDefault(logger)
}
