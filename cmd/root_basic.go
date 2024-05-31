//go:build basic && !manager

package cmd

func init() {
	rootCmd.AddCommand(appCmd)
	rootCmd.AddCommand(removeCacheCmd)
	rootCmd.AddCommand(versionsCmd)
}
