//go:build basic || manager

package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
)

var rootCmd = &cobra.Command{
	Use:   "Go Fleet Manager",
	Short: "Utilities for manage the fleet.",
	Long: `Utilities for manage the fleet.
	It has several commands to manage the different components of the fleet of Microservices`,
}

func Execute() {

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
