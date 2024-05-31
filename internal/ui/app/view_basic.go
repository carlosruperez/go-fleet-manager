//go:build basic && !manager

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
	default:
		viewString = lipgloss.JoinVertical(lipgloss.Left, m.menuModel.View(), m.help.View(m.keys))
	}

	return viewString
}
