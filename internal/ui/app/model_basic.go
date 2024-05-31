//go:build basic && !manager

package app

import (
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/go-fleet-manager/internal/ui"
	uiCommon "github.com/go-fleet-manager/internal/ui/common"
	"github.com/go-fleet-manager/internal/ui/menu"
)

type appModel struct {
	keys              uiCommon.KeyMap
	help              help.Model
	focusedView       string
	menuModel         menu.MenuModel
	versionsModel     ui.VersionsModel
	removeCacheModel  ui.RemoveCacheModel
	viewport          viewport.Model
	terminalSizeReady bool
}

func (m *appModel) resetAppModel() {
	switch m.focusedView {
	case string(menu.VersionsOption):
		m.versionsModel = ui.VersionsModel{}
	case string(menu.RemoveCacheOption):
		m.removeCacheModel = ui.RemoveCacheModel{}
	}
	m.focusedView = "menu"
}

func (m appModel) forwardMsgToView(msg tea.Msg, view string) (tea.Model, tea.Cmd) {
	switch view {
	case "menu":
		newModel, cmd := m.menuModel.Update(msg)
		m.menuModel = newModel.(menu.MenuModel)
		return m, cmd
	case string(menu.VersionsOption):
		if m.versionsModel.Equals(&ui.VersionsModel{}) {
			m.versionsModel = ui.NewVersionsModel()
		}
		newModel, cmd := m.versionsModel.Update(msg)
		m.versionsModel = newModel.(ui.VersionsModel)
		return m, cmd
	case string(menu.RemoveCacheOption):
		if (m.removeCacheModel.Equals(&ui.RemoveCacheModel{})) {
			m.removeCacheModel = ui.NewRemoveCacheModel()
		}
		newModel, cmd := m.removeCacheModel.Update(msg)
		m.removeCacheModel = newModel.(ui.RemoveCacheModel)
		return m, cmd
	}
	return m, nil

}
