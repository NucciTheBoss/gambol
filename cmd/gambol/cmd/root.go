package cmd

import (
	"errors"
	"log/slog"
	"os"
	"path"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	gambol "github.com/nuccitheboss/gambol/internal/common"
	"github.com/nuccitheboss/gambol/internal/storage"
)

const rootLongHelp = `Description:
  Playthrough complex integration tests easily

  Gambol runs integrations tests for distributed systems where
  one CI runner is not enough to validate the functionality of an application.
`

var verbose bool
var rootCmd = &cobra.Command{
	Use:     "gambol",
	Long:    rootLongHelp,
	Version: "0.1.0",
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		gambol.ConfigureLogger(verbose)
		err := storage.Init()
		if err != nil {
			slog.Error(err.Error())
			os.Exit(1)
		}
	},
}

func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "set log level to verbose")
	rootCmd.AddCommand(runCmd)
}

// Initialize gambol configuration.
func initConfig() {
	// Set `gambol` config file name and type.
	viper.SetConfigName("gambol")
	viper.SetConfigType("yaml")

	home, err := os.UserHomeDir()
	if err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}

	// Default config paths for `gambol`.
	viper.AddConfigPath(".")
	viper.AddConfigPath(path.Join(home, ".config", "gambol"))
	viper.AddConfigPath(path.Join("etc", "gambol"))

	// Set configuration defaults.
	viper.SetDefault("storage", path.Join(home, ".gambol", "storage"))

	// Read in gambol config.
	err = viper.ReadInConfig()
	if err != nil {
		if errors.As(err, &viper.ConfigFileNotFoundError{}) {
			slog.Debug("no config file found, using preset defaults")
		} else {
			slog.Error(err.Error())
			os.Exit(1)
		}
	}
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}
}
