package menu

type MenuModel struct {
	Choices        []menuChoice
	Cursor         int
	selectedChoice int
}

type menuChoice struct {
	Title    MenuOption
	selected bool
}

func (mc menuChoice) String() string {
	return string(mc.Title)
}

type MenuOption string

const (
	VersionsOption    MenuOption = "Versions"
	RemoveCacheOption MenuOption = "Remove Cache"
)
