package createMainPR

import (
	"github.com/charmbracelet/lipgloss"
)

func (m CreateMainPRModel) View() string {
	if m.focusedView == branchesOption {
		if len(m.branches.Items()) > 0 {
			return lipgloss.JoinVertical(lipgloss.Left, m.branches.View(), m.help.View(m.keys), m.errors.View())
		} else {
			return lipgloss.JoinVertical(lipgloss.Left, "🕦 loading...", m.errors.View())
		}
	}
	return lipgloss.JoinVertical(lipgloss.Left, m.CustomTableModel.View(), m.help.View(m.keys), m.errors.View())
}
