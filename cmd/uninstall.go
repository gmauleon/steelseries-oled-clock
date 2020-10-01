package cmd

import (
	"log"

	"github.com/spf13/cobra"
)

var uninstallCmd = &cobra.Command{
	Use:   "uninstall",
	Short: "Uninstall the service",
	Long:  `A service that will display the date and time on your SteelSeries devices OLED display`,
	Run: func(cmd *cobra.Command, args []string) {

		gsc, svc := GetObjects()

		if err := gsc.Unregister(); err != nil {
			log.Fatal(err)
		}

		if err := svc.Uninstall(); err != nil {
			log.Fatal(err)
		}

	},
}

func init() {
	rootCmd.AddCommand(uninstallCmd)
}
