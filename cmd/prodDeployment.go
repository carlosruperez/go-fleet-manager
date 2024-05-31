package cmd

import (
	"github.com/go-fleet-manager/internal/ui"
	"github.com/spf13/cobra"
)

var prodDeploymentCmd = &cobra.Command{
	Use:   "prodDeployment",
	Short: "Deploy to PRO.",
	Long:  `Deploy to PRO given a Release`,
	Run: func(cmd *cobra.Command, args []string) {
		ui.ProdDeploy()
	},
}
