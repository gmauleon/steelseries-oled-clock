package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "steelseries-oled-clock",
	Short: "SteelSeries OLED Clock",
	Long:  `A service that will display the date and time on your SteelSeries devices OLED display`,
}

// Execute the requested command
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {

}
