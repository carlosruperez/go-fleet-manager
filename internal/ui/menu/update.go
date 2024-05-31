package menu

import (
	tea "github.com/charmbracelet/bubbletea"
)

func (m MenuModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "up", "k":
			if m.selectedChoice != -1 {
				break
			}
			if m.Cursor > 0 {
				m.Cursor--
			} else {
				m.Cursor = len(m.Choices) - 1
			}
		case "down", "j":
			if m.selectedChoice != -1 {
				break
			}
			if m.Cursor < len(m.Choices)-1 {
				m.Cursor++
			} else {
				m.Cursor = 0
			}
		case "enter":
			for i := range m.Choices {
				if i == m.Cursor {
					m.Choices[i].selected = true
					m.selectedChoice = i
				} else {
					m.Choices[i].selected = false
				}

			}
			return m, nil
		case "esc":
			m.selectedChoice = -1
		}
	}
	return m, nil
}
