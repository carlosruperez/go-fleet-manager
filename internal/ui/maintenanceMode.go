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
	uiCommon "github.com/go-fleet-manager/internal/ui/common"
	uiCustomTable "github.com/go-fleet-manager/internal/ui/customTable"
	uiErrorsModel "github.com/go-fleet-manager/internal/ui/errors"
	"github.com/go-fleet-manager/internal/ui/styles"
	"github.com/go-fleet-manager/internal/usecases"
)

type MaintenanceModeModel struct {
	uiCustomTable.CustomTableModel
	keys   maintenanceModeKeyMap
	help   help.Model
	errors uiErrorsModel.ErrorsModel
	items  []maintenanceModeItem
}

type environmentType string

const (
	Dev     environmentType = "dev"
	Preprod environmentType = "preprod"
	Prod    environmentType = "prod"
)

type maintenanceModeItem struct {
	environment environmentType
	result      string
}

type maintenanceModeKeyMap struct {
	uiCommon.KeyMap
	Enable  key.Binding
	Disable key.Binding
}

func (k maintenanceModeKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Up, k.Down, k.Enable, k.Disable, k.Esc, k.Quit}
}

var maintenanceModeKeys = maintenanceModeKeyMap{
	KeyMap: uiCommon.Keys,
	Enable: key.NewBinding(
		key.WithKeys("e"),
		key.WithHelp(styles.KeyBindingStyle.Render("e"), "Enable maintenance"),
	),
	Disable: key.NewBinding(
		key.WithKeys("d"),
		key.WithHelp(styles.KeyBindingStyle.Render("d"), "Disable maintenance"),
	),
}

func (model *MaintenanceModeModel) Equals(other *MaintenanceModeModel) bool {
	return reflect.DeepEqual(model.items, other.items)
}

func (m *MaintenanceModeModel) updateTable() {
	values := []table.Row{}
	for _, item := range m.items {
		values = append(values, table.Row{string(item.environment), item.result})
	}
	m.SetRows(values)
}

func (m *MaintenanceModeModel) resetResults() {
	for i := range m.items {
		m.items[i].result = ""
	}
	m.updateTable()
}

func NewMaintenanceModeModel() MaintenanceModeModel {

	items := []maintenanceModeItem{
		{environment: Dev, result: ""},
		{environment: Preprod, result: ""},
		{environment: Prod, result: ""},
	}

	columns := []table.Column{
		{Title: "environment", Width: 50},
		{Title: "Result", Width: 15},
	}

	values := []table.Row{}
	for _, item := range items {
		values = append(values, table.Row{string(item.environment), item.result})
	}
	customTableModel := uiCustomTable.NewCustomTableModel(columns, values)

	m := MaintenanceModeModel{
		keys:             maintenanceModeKeys,
		help:             help.New(),
		errors:           uiErrorsModel.NewErrorsModel(),
		items:            items,
		CustomTableModel: customTableModel,
	}

	return m
}

type maintenanceModeMsg struct {
	index       int
	environment environmentType
	result      string
	err         error
}

func runMaintenanceMode(index int, environment environmentType, action usecases.MaintenanceAction) tea.Cmd {
	ctx, cancel := context.WithCancel(context.Background()) // TODO: use context to get feedback
	go func() tea.Msg {
		<-ctx.Done()
		// The context was cancelled, stop the update process
		// and return a special message
		return maintenanceModeMsg{index, environment, "❌" + "Cancelled", nil}
	}()
	return func() tea.Msg {
		mapCommonEnvironment := map[environmentType]common.Environment{
			Dev:     common.Development,
			Preprod: common.Staging,
			Prod:    common.Production,
		}
		commonEnvironment := mapCommonEnvironment[environment]
		_, _, err := usecases.MaintenanceMode(action, commonEnvironment, ctx)
		if err != nil {
			cancel()
			return maintenanceModeMsg{index, environment, "❌", err}
		}
		return maintenanceModeMsg{index, environment, "🚀", nil}
	}
}

func (m MaintenanceModeModel) Init() tea.Cmd {
	return nil
}

func (m MaintenanceModeModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {

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
		case key.Matches(msg, m.keys.Enable):
			return m, runMaintenanceMode(m.Cursor(), m.items[m.Cursor()].environment, usecases.EnableMaintenance)
		case key.Matches(msg, m.keys.Disable):
			return m, runMaintenanceMode(m.Cursor(), m.items[m.Cursor()].environment, usecases.DisableMaintenance)
		case key.Matches(msg, m.keys.Esc):
			m.resetResults()
			return m, nil
		}
	case maintenanceModeMsg:
		m.items[msg.index].result = msg.result
		if msg.err != nil {
			m.errors.AddError(msg.err)
		}
		m.updateTable()
	}
	return m, nil
}

func (m MaintenanceModeModel) View() string {
	return lipgloss.JoinVertical(lipgloss.Left, m.CustomTableModel.View(), m.help.View(m.keys), m.errors.View())
}

func MaintenanceMode() {
	maintenanceModeModel := NewMaintenanceModeModel()
	if _, err := tea.NewProgram(maintenanceModeModel).Run(); err != nil {
		fmt.Printf("Uh oh, there was an error: %v\n", err)
		os.Exit(1)
	}
}
