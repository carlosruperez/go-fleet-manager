package cmd

import (
	"github.com/go-fleet-manager/internal/ui/app"
	"github.com/spf13/cobra"
)

var appCmd = &cobra.Command{
	Use:   "app",
	Short: "The App manage our fleet.",
	Long:  "The App to manage our fleet.",
	Run: func(cmd *cobra.Command, args []string) {
		app.RunApp()
	},
}
