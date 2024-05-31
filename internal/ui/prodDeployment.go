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

type ProdDeploymentModel struct {
	uiCustomTable.CustomTableModel
	keys   prodDeploymentKeyMap
	help   help.Model
	errors uiErrorsModel.ErrorsModel
	items  []prodDeploymentItem
}

type prodDeploymentItem struct {
	repository repository.Repository
	result     string
}

type prodDeploymentKeyMap struct {
	uiCommon.KeyMap
	Enter key.Binding
}

func (k prodDeploymentKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Up, k.Down, k.Enter, k.Esc, k.Quit}
}

var prodDeploymentKeys = prodDeploymentKeyMap{
	KeyMap: uiCommon.Keys,
	Enter: key.NewBinding(
		key.WithKeys("enter"),
		key.WithHelp(styles.KeyBindingStyle.Render("enter"), "Deploy to PROD"),
	),
}

func (model *ProdDeploymentModel) Equals(other *ProdDeploymentModel) bool {
	return reflect.DeepEqual(model.items, other.items)
}

func (m *ProdDeploymentModel) updateTable() {
	values := []table.Row{}
	for i, item := range m.items {
		values = append(values, table.Row{fmt.Sprint(i + 1), item.repository.Name, item.result})
	}
	m.SetRows(values)
}

func (m *ProdDeploymentModel) resetResults() {
	for i := range m.items {
		m.items[i].result = ""
	}
	m.updateTable()
}

func NewProdDeploymentModel() ProdDeploymentModel {

	items := []prodDeploymentItem{}

	repositories := usecases.GetRepositories()

	for _, repo := range repositories {
		if repo.Type != common.Microservice {
			continue
		}
		items = append(items, prodDeploymentItem{repository: repo})
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

	m := ProdDeploymentModel{
		keys:             prodDeploymentKeys,
		help:             help.New(),
		errors:           uiErrorsModel.NewErrorsModel(),
		items:            items,
		CustomTableModel: customTableModel,
	}

	return m
}

type prodDeploymentMsg struct {
	index  int
	repo   repository.Repository
	result string
	err    error
}

func prodDeploy(index int, repo repository.Repository) tea.Cmd {
	ctx, cancel := context.WithCancel(context.Background()) // TODO: use context to get feedback
	go func() tea.Msg {
		<-ctx.Done()
		// The context was cancelled, stop the update process
		// and return a special message
		return prodDeploymentMsg{index, repo, "❌" + "Cancelled", nil}
	}()
	return func() tea.Msg {
		version, err := usecases.GetVersion(repo, common.Staging)
		if err != nil {
			cancel()
			return prodDeploymentMsg{index, repo, "❌", err}
		}
		_, _, err = usecases.ProdDeploy(repo, version, ctx)
		if err != nil {
			cancel()
			return prodDeploymentMsg{index, repo, "❌", err}
		}
		return prodDeploymentMsg{index, repo, "🚀", nil}
	}
}

func (m ProdDeploymentModel) Init() tea.Cmd {
	return nil
}

func (m ProdDeploymentModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {

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
			return m, prodDeploy(m.Cursor(), m.items[m.Cursor()].repository)
		case key.Matches(msg, m.keys.Esc):
			m.resetResults()
			return m, nil
		}
	case prodDeploymentMsg:
		m.items[msg.index].result = msg.result
		if msg.err != nil {
			m.errors.AddError(msg.err)
		}
		m.updateTable()
	}
	return m, nil
}

func (m ProdDeploymentModel) View() string {
	return lipgloss.JoinVertical(lipgloss.Left, m.CustomTableModel.View(), m.help.View(m.keys), m.errors.View())
}

func ProdDeploy() {
	prodDeploymentModel := NewProdDeploymentModel()
	if _, err := tea.NewProgram(prodDeploymentModel).Run(); err != nil {
		fmt.Printf("Uh oh, there was an error: %v\n", err)
		os.Exit(1)
	}
}
