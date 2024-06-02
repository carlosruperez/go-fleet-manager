package ui

import (
	"context"
	"fmt"
	"os"
	"reflect"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/go-fleet-manager/internal/common"
	"github.com/go-fleet-manager/internal/repository"
	uiCommon "github.com/go-fleet-manager/internal/ui/common"
	uiCustomTable "github.com/go-fleet-manager/internal/ui/customTable"
	uiErrorsModel "github.com/go-fleet-manager/internal/ui/errors"
	"github.com/go-fleet-manager/internal/ui/styles"
	"github.com/go-fleet-manager/internal/usecases"
)

type DeploymentModel struct {
	uiCustomTable.CustomTableModel
	keys   deploymentKeyMap
	help   help.Model
	errors uiErrorsModel.ErrorsModel
	items  []deploymentItem
}

type deploymentItem struct {
	repository repository.Repository
	result     string
}

type deploymentKeyMap struct {
	uiCommon.KeyMap
	Enter key.Binding
}

func (k deploymentKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Up, k.Down, k.Enter, k.Esc, k.Quit}
}

var deploymentKeys = deploymentKeyMap{
	KeyMap: uiCommon.Keys,
	Enter: key.NewBinding(
		key.WithKeys("enter"),
		key.WithHelp(styles.KeyBindingStyle.Render("enter"), "Deploy to PROD"),
	),
}

func (model *DeploymentModel) Equals(other *DeploymentModel) bool {
	return reflect.DeepEqual(model.items, other.items)
}

func (m *DeploymentModel) updateTable() {
	values := []table.Row{}
	for i, item := range m.items {
		values = append(values, table.Row{fmt.Sprint(i + 1), item.repository.Name, item.result})
	}
	m.SetRows(values)
}

func (m *DeploymentModel) resetResults() {
	for i := range m.items {
		m.items[i].result = ""
	}
	m.updateTable()
}

func NewDeploymentModel() DeploymentModel {

	items := []deploymentItem{}

	repositories := usecases.GetRepositories()

	for _, repo := range repositories {
		if repo.Type != common.Microservice {
			continue
		}
		items = append(items, deploymentItem{repository: repo})
	}

	columns := []table.Column{
		{Title: "n°", Width: 5},
		{Title: "Repository", Width: 50},
		{Title: "Result", Width: 15},
	}

	values := []table.Row{}
	for i, item := range items {
		values = append(values, table.Row{fmt.Sprint(i + 1), item.repository.Name, item.result})
	}
	customTableModel := uiCustomTable.NewCustomTableModel(columns, values)

	m := DeploymentModel{
		keys:             deploymentKeys,
		help:             help.New(),
		errors:           uiErrorsModel.NewErrorsModel(),
		items:            items,
		CustomTableModel: customTableModel,
	}

	return m
}

type deploymentMsg struct {
	index  int
	repo   repository.Repository
	result string
	err    error
}

func deploy(index int, repo repository.Repository) tea.Cmd {
	ctx, cancel := context.WithCancel(context.Background()) // TODO: use context to get feedback
	go func() tea.Msg {
		<-ctx.Done()
		// The context was cancelled, stop the update process
		// and return a special message
		return deploymentMsg{index, repo, "❌" + "Cancelled", nil}
	}()
	return func() tea.Msg {
		version, err := usecases.GetVersion(repo, common.Staging)
		if err != nil {
			cancel()
			return deploymentMsg{index, repo, "❌", err}
		}
		_, _, err = usecases.Deploy(repo, version, ctx)
		if err != nil {
			cancel()
			return deploymentMsg{index, repo, "❌", err}
		}
		return deploymentMsg{index, repo, "🚀", nil}
	}
}

func (m DeploymentModel) Init() tea.Cmd {
	return nil
}

func (m DeploymentModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {

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
		case key.Matches(msg, m.keys.Enter):
			return m, deploy(m.Cursor(), m.items[m.Cursor()].repository)
		case key.Matches(msg, m.keys.Esc):
			m.resetResults()
			return m, nil
		}
	case deploymentMsg:
		m.items[msg.index].result = msg.result
		if msg.err != nil {
			m.errors.AddError(msg.err)
		}
		m.updateTable()
	}
	return m, nil
}

func (m DeploymentModel) View() string {
	return lipgloss.JoinVertical(lipgloss.Left, m.CustomTableModel.View(), m.help.View(m.keys), m.errors.View())
}

func Deploy() {
	deploymentModel := NewDeploymentModel()
	if _, err := tea.NewProgram(deploymentModel).Run(); err != nil {
		fmt.Printf("Uh oh, there was an error: %v\n", err)
		os.Exit(1)
	}
}
