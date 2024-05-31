package cmd

import (
	"github.com/go-fleet-manager/internal/ui"
	"github.com/spf13/cobra"
)

var versionsTomlVsDeployCmd = &cobra.Command{
	Use:   "versions-toml-vs-deploy",
	Short: "Get versions of the different applications of .toml and deployed.",
	Long:  `Get versions of the different applications of .toml and deployed`,
	Run: func(cmd *cobra.Command, args []string) {
		ui.CheckVersions()
	},
}
