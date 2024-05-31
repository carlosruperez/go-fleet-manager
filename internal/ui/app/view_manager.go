//go:build manager && !basic

package app

import (
	"github.com/charmbracelet/lipgloss"
	"github.com/go-fleet-manager/internal/ui/menu"
)

func (m appModel) view() string {

	var viewString string

	switch {
	case m.focusedView == string(menu.VersionsOption):
		viewString = lipgloss.JoinHorizontal(lipgloss.Left, m.menuModel.View(), m.versionsModel.View())
	case m.focusedView == string(menu.RemoveCacheOption):
		viewString = lipgloss.JoinHorizontal(lipgloss.Left, m.menuModel.View(), m.removeCacheModel.View())
	case m.focusedView == string(menu.MaintenanceOption):
		viewString = lipgloss.JoinHorizontal(lipgloss.Left, m.menuModel.View(), m.maintenanceMode.View())
	case m.focusedView == string(menu.DevMinorOption):
		viewString = lipgloss.JoinHorizontal(lipgloss.Left, m.menuModel.View(), m.devMinorModel.View())
	case m.focusedView == string(menu.PRODDeploymentOption):
		viewString = lipgloss.JoinHorizontal(lipgloss.Left, m.menuModel.View(), m.prodDeploymentModel.View())
	case m.focusedView == string(menu.CreatePRsToMainOption):
		viewString = lipgloss.JoinHorizontal(lipgloss.Left, m.menuModel.View(), m.createMainPRModel.View())
	case m.focusedView == string(menu.VersionsTomlvsDeployOption):
		viewString = lipgloss.JoinHorizontal(lipgloss.Left, m.menuModel.View(), m.versionsTomlVsDeployModel.View())
	default:
		viewString = lipgloss.JoinVertical(lipgloss.Left, m.menuModel.View(), m.help.View(m.keys))
	}

	return viewString
}
