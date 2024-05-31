//go:build (manager && !basic) || (basic && !manager)

package app

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const (
	minTerminalWidth  = 160
	minTerminalHeight = 35
)

func (m *appModel) syncTerminal(msg tea.Msg) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		headerHeight := lipgloss.Height(m.headerView())

		if !m.terminalSizeReady {
			m.viewport = viewport.New(msg.Width, msg.Height)
			m.viewport.YPosition = headerHeight + 1
			m.terminalSizeReady = true
		} else {
			m.viewport.Width = msg.Width
			m.viewport.Height = msg.Height
		}
	}
}

func (m *appModel) headerView(titles ...string) string {
	var renderedTitles string
	for _, t := range titles {
		renderedTitles += t
	}
	line := strings.Repeat("─", max(0, m.viewport.Width-79))
	titles = append(titles, line)
	return lipgloss.JoinHorizontal(lipgloss.Center, titles...)
}

func (m appModel) View() string {

	if !m.terminalSizeReady {
		return "Setting up..."
	}

	if m.viewport.Width < minTerminalWidth || m.viewport.Height < minTerminalHeight {
		return fmt.Sprintf("Terminal window is too small. Please resize to at least %dx%d.", minTerminalWidth, minTerminalHeight)
	}

	viewString := m.view()

	return viewString
}
