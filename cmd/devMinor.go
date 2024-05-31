package cmd

import (
	"github.com/go-fleet-manager/internal/ui"
	"github.com/spf13/cobra"
)

var devMinorCmd = &cobra.Command{
	Use:   "dev-minor",
	Short: "Update the minor version of the different applications.",
	Long:  `Update the minor version of the different applications.`,
	Run: func(cmd *cobra.Command, args []string) {
		ui.DevMinor()
	},
}
