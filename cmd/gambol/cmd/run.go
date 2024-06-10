package cmd

import (
	"log/slog"
	"os"

	"github.com/spf13/cobra"

	gambol "github.com/nuccitheboss/gambol/internal/common"
)

const runShortHelp = "Run gambol playthroughs"
const runLongHelp = `Description:
  Run a gambol playthrough
`
const examples = `  gambol run spec.yaml
      Run gambol playthrough specificed in spec.yaml
`

var runCmd = &cobra.Command{
	Use:     "run",
	Short:   runShortHelp,
	Long:    runLongHelp,
	Example: examples,
	Args:    cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		play := args[0]
		_, err := os.Stat(play)
		if err != nil {
			slog.Error(err.Error())
			os.Exit(2)
		}

		err = gambol.Run(play)
		if err != nil {
			slog.Error(err.Error())
			os.Exit(1)
		}
	},
}
