package cmd

import (
	"github.com/go-fleet-manager/internal/ui"
	"github.com/spf13/cobra"
)

var deploymentCmd = &cobra.Command{
	Use:   "deployment",
	Short: "Deploy to an environment.",
	Long:  `Deploy to an environment given a Release`,
	Run: func(cmd *cobra.Command, args []string) {
		ui.Deploy()
	},
}
