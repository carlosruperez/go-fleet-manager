//go:build manager && !basic

package cmd

func init() {
	rootCmd.AddCommand(appCmd)
	rootCmd.AddCommand(createMainPRCmd)
	rootCmd.AddCommand(devMinorCmd)
	rootCmd.AddCommand(maintenanceModeCmd)
	rootCmd.AddCommand(prodDeploymentCmd)
	rootCmd.AddCommand(removeCacheCmd)
	rootCmd.AddCommand(versionsCmd)
}
