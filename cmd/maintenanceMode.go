package cmd

import (
	"github.com/go-fleet-manager/internal/ui"
	"github.com/spf13/cobra"
)

var maintenanceModeCmd = &cobra.Command{
	Use:   "maintenanceMode",
	Short: "Enable/Disable maintenance mode.",
	Long:  `Enable/Disable maintenance mode.`,
	Run: func(cmd *cobra.Command, args []string) {
		ui.MaintenanceMode()
	},
}
