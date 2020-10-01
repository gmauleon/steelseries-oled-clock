package cmd

import (
	"fmt"
	"log"
	"os"

	"github.com/gmauleon/steelseries-oled-clock/pkg/clock"
	"github.com/kardianos/service"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "steelseries-oled-clock",
	Short: "SteelSeries OLED Clock",
	Long:  `A service that will display the date and time on your SteelSeries devices OLED display`,
	Run: func(cmd *cobra.Command, args []string) {
		gsc, svc := GetObjects()
		if err := svc.Run(); err != nil {
			gsc.Logger.Error(err)
		}
	},
}

// GetObjects return a clock and the corresponding service
func GetObjects() (gsc *clock.GameSenseClock, svc service.Service) {
	svcConfig := &service.Config{
		Name:        "SteelSeriesOLEDClock",
		DisplayName: "SteelSeries OLED Clock",
		Description: "This service communicate with GameSense to display the current date and time.",
	}

	gsc = clock.NewGameSenseClock()

	svc, err := service.New(gsc, svcConfig)
	if err != nil {
		log.Fatal(err)
	}

	logger, err := svc.Logger(nil)
	if err != nil {
		log.Fatal(err)
	}

	gsc.Init(logger)
	return
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
