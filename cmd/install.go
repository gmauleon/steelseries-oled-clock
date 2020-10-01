package cmd

import (
	"log"

	"github.com/spf13/cobra"
)

var installCmd = &cobra.Command{
	Use:   "install",
	Short: "Install the service",
	Run: func(cmd *cobra.Command, args []string) {

		gsc, svc := GetObjects()

		if err := svc.Install(); err != nil {
			log.Fatal(err)
		}

		if err := gsc.Register(); err != nil {
			log.Fatal(err)
		}

	},
}

func init() {
	rootCmd.AddCommand(installCmd)
}
