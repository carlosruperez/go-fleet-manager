package cmd

import (
	"github.com/go-fleet-manager/internal/ui"
	"github.com/spf13/cobra"
)

var versionsCmd = &cobra.Command{
	Use:   "versions",
	Short: "Get versions of the different applications.",
	Long:  `Get versions of the different applications by environment.`,
	Run: func(cmd *cobra.Command, args []string) {
		ui.GetVersions()
	},
}
