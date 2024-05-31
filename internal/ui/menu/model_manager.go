//go:build manager && !basic

package menu

const (
	MaintenanceOption          MenuOption = "Maintenance Mode"
	MigrationsOption           MenuOption = "Migrations"
	DevMinorOption             MenuOption = "DevMinor"
	CreatePRsToMainOption      MenuOption = "Create PRs to main"
	PRODDeploymentOption       MenuOption = "PROD Deployment"
	PRODDeploymentPlanOption   MenuOption = "PROD Deployment plan"
	VersionsTomlvsDeployOption MenuOption = "Versions Toml vs Deploy"
)

func NewMenuModel() MenuModel {
	choices := []menuChoice{
		{Title: VersionsOption},
		{Title: RemoveCacheOption},
		{Title: MaintenanceOption},
		{Title: MigrationsOption},
		{Title: DevMinorOption},
		{Title: CreatePRsToMainOption},
		{Title: PRODDeploymentOption},
		{Title: PRODDeploymentPlanOption},
		{Title: VersionsTomlvsDeployOption},
	}
	return MenuModel{Choices: choices, selectedChoice: -1}
}
