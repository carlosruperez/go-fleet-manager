//go:build manager && !basic

package app

import (
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/go-fleet-manager/internal/ui"
	uiCommon "github.com/go-fleet-manager/internal/ui/common"
	"github.com/go-fleet-manager/internal/ui/createMainPR"
	"github.com/go-fleet-manager/internal/ui/menu"
)

type appModel struct {
	keys                      uiCommon.KeyMap
	help                      help.Model
	focusedView               string
	menuModel                 menu.MenuModel
	versionsModel             ui.VersionsModel
	removeCacheModel          ui.RemoveCacheModel
	maintenanceMode           ui.MaintenanceModeModel
	devMinorModel             ui.DevMinorModel
	prodDeploymentModel       ui.ProdDeploymentModel
	createMainPRModel         createMainPR.CreateMainPRModel
	versionsTomlVsDeployModel ui.VersionsTomlVsDeployModel
	viewport                  viewport.Model
	terminalSizeReady         bool
}

func (m *appModel) resetAppModel() {
	switch m.focusedView {
	case string(menu.VersionsOption):
		m.versionsModel = ui.VersionsModel{}
	case string(menu.RemoveCacheOption):
		m.removeCacheModel = ui.RemoveCacheModel{}
	case string(menu.MaintenanceOption):
		m.maintenanceMode = ui.MaintenanceModeModel{}
	case string(menu.DevMinorOption):
		m.devMinorModel = ui.DevMinorModel{}
	case string(menu.PRODDeploymentOption):
		m.prodDeploymentModel = ui.ProdDeploymentModel{}
	case string(menu.CreatePRsToMainOption):
		m.createMainPRModel = createMainPR.CreateMainPRModel{}
	case string(menu.VersionsTomlvsDeployOption):
		m.versionsTomlVsDeployModel = ui.VersionsTomlVsDeployModel{}
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
	case string(menu.MaintenanceOption):
		if (m.maintenanceMode.Equals(&ui.MaintenanceModeModel{})) {
			m.maintenanceMode = ui.NewMaintenanceModeModel()
		}
		newModel, cmd := m.maintenanceMode.Update(msg)
		m.maintenanceMode = newModel.(ui.MaintenanceModeModel)
		return m, cmd
	case string(menu.DevMinorOption):
		if (m.devMinorModel.Equals(&ui.DevMinorModel{})) {
			m.devMinorModel = ui.NewDevMinorModel()
		}
		newModel, cmd := m.devMinorModel.Update(msg)
		m.devMinorModel = newModel.(ui.DevMinorModel)
		return m, cmd
	case string(menu.PRODDeploymentOption):
		if (m.prodDeploymentModel.Equals(&ui.ProdDeploymentModel{})) {
			m.prodDeploymentModel = ui.NewProdDeploymentModel()
		}
		newModel, cmd := m.prodDeploymentModel.Update(msg)
		m.prodDeploymentModel = newModel.(ui.ProdDeploymentModel)
		return m, cmd
	case string(menu.CreatePRsToMainOption):
		if (m.createMainPRModel.Equals(&createMainPR.CreateMainPRModel{})) {
			m.createMainPRModel = createMainPR.NewCreateMainPRModel()
		}
		newModel, cmd := m.createMainPRModel.Update(msg)
		m.createMainPRModel = newModel.(createMainPR.CreateMainPRModel)
		return m, cmd
	case string(menu.VersionsTomlvsDeployOption):
		if (m.versionsTomlVsDeployModel.Equals(&ui.VersionsTomlVsDeployModel{})) {
			m.versionsTomlVsDeployModel = ui.NewVersionsTomlVsDeployModel()
		}
		newModel, cmd := m.versionsTomlVsDeployModel.Update(msg)
		m.versionsTomlVsDeployModel = newModel.(ui.VersionsTomlVsDeployModel)
		return m, cmd
	}
	return m, nil

}
