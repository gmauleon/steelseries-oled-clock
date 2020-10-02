package cmd

import (
	"github.com/gmauleon/steelseries-oled-clock/pkg/clock"
	"github.com/spf13/cobra"
)

var uninstallCmd = &cobra.Command{
	Use:   "uninstall",
	Short: "Uninstall the service",
	Long:  `A service that will display the date and time on your SteelSeries devices OLED display`,
	Run: func(cmd *cobra.Command, args []string) {
		gs := clock.NewGameSenseClockService()
		gs.UninstallService()
	},
}

func init() {
	rootCmd.AddCommand(uninstallCmd)
}
