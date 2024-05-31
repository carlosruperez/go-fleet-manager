package createMainPR

import (
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

func (m CreateMainPRModel) updateBranchesList(msg tea.Msg) (CreateMainPRModel, tea.Cmd) {
	switch msg := msg.(type) {
	case selectRepositoryBranchMsg:
		var branchItems []list.Item

		for _, item := range msg.branches {
			branchItems = append(branchItems, branchItem{name: item})
		}
		m.branches.SetItems(branchItems)
		if msg.err != nil {
			m.errors.AddError(msg.err)
		}
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.keys.Quit):
			return m, tea.Quit
		case key.Matches(msg, m.keys.Enter):
			branchName := m.branches.SelectedItem().(branchItem).name
			m.branches.SetItems(nil)
			m.focusedView = tableOption
			m.items[m.selectedItem].selectedBranch = branchName
			m.updateTable()
			return m, nil
		case key.Matches(msg, m.keys.Esc):
			m.branches.SetItems(nil)
			m.focusedView = tableOption
			return m, nil
		}
	}
	newListModel, cmd := m.branches.Update(msg)
	m.branches = newListModel
	return m, cmd
}

func (m CreateMainPRModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if m.focusedView == branchesOption {
		return m.updateBranchesList(msg)
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.keys.Quit):
			return m, tea.Quit
		case key.Matches(msg, m.keys.Up):
			if m.Cursor() > 0 {
				m.SetCursor(m.Cursor() - 1)
			}
		case key.Matches(msg, m.keys.Down):
			if m.Cursor() < len(m.items)-1 {
				m.SetCursor(m.Cursor() + 1)
			}
		case key.Matches(msg, m.keys.Select):
			m.focusedView = branchesOption
			m.selectedItem = m.Cursor()
			return m, initSelectRepositoryBranch(m.items[m.selectedItem].repository, m.selectedItem)
		case key.Matches(msg, m.keys.Run):
			return m, createMainPR(m.Cursor(), m.items[m.Cursor()].repository, m.items[m.Cursor()].selectedBranch)
		case key.Matches(msg, m.keys.Enter):
			cmds := make([]tea.Cmd, 0)

			for index, item := range m.items {
				if item.selectedBranch != "" {
					cmds = append(cmds, createMainPR(index, item.repository, item.selectedBranch))
				}
			}
			return m, tea.Batch(cmds...)
		case key.Matches(msg, m.keys.Esc):
			m.resetResults()
			return m, nil
		}
	case createMainPRMsg:
		m.items[msg.index].result = msg.result
		if msg.err != nil {
			m.errors.AddError(msg.err)
		}
		m.updateTable()
	}
	return m, nil
}
