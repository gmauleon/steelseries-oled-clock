package cmd

import (
	"github.com/gmauleon/steelseries-oled-clock/pkg/clock"
	"github.com/spf13/cobra"
)

var installCmd = &cobra.Command{
	Use:   "install",
	Short: "Install the service",
	Run: func(cmd *cobra.Command, args []string) {
		gs := clock.NewGameSenseClockService()
		gs.InstallService()
	},
}

func init() {
	rootCmd.AddCommand(installCmd)
}
