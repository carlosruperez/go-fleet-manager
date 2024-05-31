package cmd

import (
	"github.com/go-fleet-manager/internal/ui/createMainPR"
	"github.com/spf13/cobra"
)

var createMainPRCmd = &cobra.Command{
	Use:   "createMainPR",
	Short: "Create a PR to main.",
	Long:  `Create a PR to main.`,
	Run: func(cmd *cobra.Command, args []string) {
		createMainPR.CreateMainPR()
	},
}
