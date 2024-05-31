//go:build basic || manager

package app

import (
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/go-fleet-manager/internal/ui/menu"
)

func (m appModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	m.syncTerminal(msg)

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.keys.Esc):
			m.resetAppModel()
		case key.Matches(msg, m.keys.Enter):
			if m.focusedView != "menu" {
				return m.forwardMsgToView(msg, m.focusedView)
			}
			m.focusedView = m.menuModel.Choices[m.menuModel.Cursor].String()
			switch m.menuModel.Choices[m.menuModel.Cursor].Title {
			case menu.VersionsOption:
				model, _ := m.forwardMsgToView(msg, m.focusedView)
				return model, m.versionsModel.Init()
			}
			return m.forwardMsgToView(nil, m.focusedView)
		case key.Matches(msg, m.keys.Help):
			m.help.ShowAll = !m.help.ShowAll
		default:
			return m.forwardMsgToView(msg, m.focusedView)
		}
	}
	return m.forwardMsgToView(msg, m.focusedView)
}
