//go:build basic || manager

package app

import (
	"github.com/charmbracelet/bubbles/help"
	uiCommon "github.com/go-fleet-manager/internal/ui/common"
	"github.com/go-fleet-manager/internal/ui/menu"
)

func newAppModel() appModel {
	return appModel{
		keys:        uiCommon.Keys,
		help:        help.New(),
		focusedView: "menu",
		menuModel:   menu.NewMenuModel(),
	}
}
