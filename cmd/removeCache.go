package cmd

import (
	"github.com/go-fleet-manager/internal/ui"
	"github.com/spf13/cobra"
)

var removeCacheCmd = &cobra.Command{
	Use:   "remove-cache",
	Short: "Remove all the stored data from the Cache.",
	Long:  `Be careful with this command. Remove all the stored data from the Cache.`,
	Run: func(cmd *cobra.Command, args []string) {
		ui.RemoveAllCache()
	},
}
