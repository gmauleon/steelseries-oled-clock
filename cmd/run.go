package cmd

import (
	"github.com/gmauleon/steelseries-oled-clock/pkg/clock"
	"github.com/spf13/cobra"
)

var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Run the service",
	Run: func(cmd *cobra.Command, args []string) {
		gs := clock.NewGameSenseClockService()
		gs.RunService()
	},
}

func init() {
	rootCmd.AddCommand(runCmd)
}
