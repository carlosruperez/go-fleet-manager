//go:build basic && !manager

package menu

func NewMenuModel() MenuModel {
	choices := []menuChoice{
		{Title: VersionsOption},
		{Title: RemoveCacheOption},
	}
	return MenuModel{Choices: choices, selectedChoice: -1}
}
