package ui

import (
	"context"
	"fmt"
	"github.com/charmbracelet/bubbles/list"
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

type viewOption int

const (
	tableOption viewOption = iota
	environmentOption
	tagsOption
)

type DeploymentModel struct {
	uiCustomTable.CustomTableModel
	keys                deploymentKeyMap
	help                help.Model
	errors              uiErrorsModel.ErrorsModel
	items               []deploymentItem
	selectedItem        int
	selectedEnvironment common.Environment
	tags                list.Model
	environments        list.Model
	focusedView         viewOption
}

type tagItem struct {
	name string
	list.Item
}

func (t tagItem) Title() string       { return t.name }
func (t tagItem) Description() string { return t.name }
func (t tagItem) FilterValue() string { return t.name }

type environmentItem struct {
	name        string
	environment common.Environment
	list.Item
}

func (t environmentItem) Title() string       { return t.name }
func (t environmentItem) Description() string { return t.name }
func (t environmentItem) FilterValue() string { return t.name }

type deploymentItem struct {
	repository  repository.Repository
	selectedTag string
	result      string
}

type deploymentKeyMap struct {
	uiCommon.KeyMap
	Enter  key.Binding
	Run    key.Binding
	Select key.Binding
}

func (k deploymentKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Up, k.Down, k.Select, k.Run, k.Enter, k.Esc, k.Quit}
}

var deploymentKeys = deploymentKeyMap{
	KeyMap: uiCommon.Keys,
	Run: key.NewBinding(
		key.WithKeys("r"),
		key.WithHelp(styles.KeyBindingStyle.Render("r"), "Run selected option"),
	),
	Select: key.NewBinding(
		key.WithKeys("s"),
		key.WithHelp(styles.KeyBindingStyle.Render("s"), "Select tag"),
	),
	Enter: key.NewBinding(
		key.WithKeys("enter"),
		key.WithHelp(styles.KeyBindingStyle.Render("enter"), "Deploy All"),
	),
}

func (model *DeploymentModel) Equals(other *DeploymentModel) bool {
	return reflect.DeepEqual(model.items, other.items)
}

func (m *DeploymentModel) updateTable() {
	values := []table.Row{}
	for i, item := range m.items {
		values = append(values, table.Row{fmt.Sprint(i + 1), item.repository.Name, item.selectedTag, item.result})
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
		{Title: "Selected Tag", Width: 25},
		{Title: "Result", Width: 15},
	}

	values := []table.Row{}
	for i, item := range items {
		values = append(values, table.Row{fmt.Sprint(i + 1), item.repository.Name, "", item.result})
	}
	customTableModel := uiCustomTable.NewCustomTableModel(columns, values)

	delegate := list.NewDefaultDelegate()
	delegate.ShowDescription = false
	newList := list.New([]list.Item{}, delegate, 50, 13)
	newList.Title = "Select a tag version to deploy:"
	newList.SetShowHelp(false)
	newList.SetShowFilter(false)
	newList.SetShowPagination(false)
	newList.SetShowStatusBar(false)

	environmentsDelegate := list.NewDefaultDelegate()
	environmentsDelegate.ShowDescription = false
	environmentItems := []list.Item{
		environmentItem{name: "dev", environment: common.Development},
		environmentItem{name: "pre", environment: common.Staging},
		environmentItem{name: "prod", environment: common.Production},
	}
	environmentsList := list.New(environmentItems, delegate, 55, 13)
	environmentsList.Title = "Select an Environment to deploy the tags version:"
	environmentsList.SetShowHelp(false)
	environmentsList.SetShowFilter(false)
	environmentsList.SetShowPagination(false)
	environmentsList.SetShowStatusBar(false)

	m := DeploymentModel{
		keys:             deploymentKeys,
		help:             help.New(),
		errors:           uiErrorsModel.NewErrorsModel(),
		items:            items,
		CustomTableModel: customTableModel,
		tags:             newList,
		environments:     environmentsList,
		focusedView:      environmentOption,
	}

	return m
}

type deploymentMsg struct {
	index  int
	repo   repository.Repository
	result string
	err    error
}

type selectRepositoryTagsMsg struct {
	selectedItem int
	tags         []string
	err          error
}

func initSelectRepositoryTags(repo repository.Repository, selectedItem int) tea.Cmd {
	ctx, cancel := context.WithCancel(context.Background()) // TODO: use context to get feedback
	go func() tea.Msg {
		<-ctx.Done()
		// The context was cancelled, stop the update process
		// and return a special message
		return selectRepositoryTagsMsg{selectedItem, []string{"❌" + "Cancelled"}, nil}
	}()
	return func() tea.Msg {
		tags, err := usecases.GetRepositoryTags(repo, ctx)
		if err != nil {
			cancel()
			return selectRepositoryTagsMsg{selectedItem, nil, err}
		}
		return selectRepositoryTagsMsg{selectedItem, tags, nil}
	}
}

func deploy(index int, repo repository.Repository, tagVersion string, environment common.Environment) tea.Cmd {
	ctx, cancel := context.WithCancel(context.Background()) // TODO: use context to get feedback
	go func() tea.Msg {
		<-ctx.Done()
		// The context was cancelled, stop the update process
		// and return a special message
		return deploymentMsg{index, repo, "❌" + "Cancelled", nil}
	}()
	return func() tea.Msg {
		_, _, err := usecases.Deploy(repo, tagVersion, environment, ctx)
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

func (m DeploymentModel) updateTagsList(msg tea.Msg) (DeploymentModel, tea.Cmd) {
	switch msg := msg.(type) {
	case selectRepositoryTagsMsg:
		var tagItems []list.Item

		for _, item := range msg.tags {
			tagItems = append(tagItems, tagItem{name: item})
		}
		m.tags.SetItems(tagItems)
		if msg.err != nil {
			m.errors.AddError(msg.err)
		}
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.keys.Quit):
			return m, tea.Quit
		case key.Matches(msg, m.keys.Enter):
			tagName := m.tags.SelectedItem().(tagItem).name
			m.tags.SetItems(nil)
			m.focusedView = tableOption
			m.items[m.selectedItem].selectedTag = tagName
			m.updateTable()
			return m, nil
		case key.Matches(msg, m.keys.Esc):
			m.tags.SetItems(nil)
			m.focusedView = tableOption
			return m, nil
		}
	}
	newListModel, cmd := m.tags.Update(msg)
	m.tags = newListModel
	return m, cmd
}

func (m DeploymentModel) updateEnvironmentList(msg tea.Msg) (DeploymentModel, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.keys.Quit):
			return m, tea.Quit
		case key.Matches(msg, m.keys.Enter):
			selectedEnvironment := m.environments.SelectedItem().(environmentItem)
			m.focusedView = tableOption
			m.selectedEnvironment = selectedEnvironment.environment
			return m, nil
		}
	}
	newListModel, cmd := m.environments.Update(msg)
	m.environments = newListModel
	return m, cmd
}

func (m DeploymentModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {

	if m.focusedView == tagsOption {
		return m.updateTagsList(msg)
	}

	if m.focusedView == environmentOption {
		return m.updateEnvironmentList(msg)
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
			m.focusedView = tagsOption
			m.selectedItem = m.Cursor()
			return m, initSelectRepositoryTags(m.items[m.selectedItem].repository, m.selectedItem)
		case key.Matches(msg, m.keys.Run):
			return m, deploy(m.Cursor(), m.items[m.Cursor()].repository, m.items[m.Cursor()].selectedTag, m.selectedEnvironment)
		case key.Matches(msg, m.keys.Enter):
			cmds := make([]tea.Cmd, 0)

			for index, item := range m.items {
				if item.selectedTag != "" {
					cmds = append(cmds, deploy(index, item.repository, item.selectedTag, m.selectedEnvironment))
				}
			}
			return m, tea.Batch(cmds...)
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
	if m.focusedView == environmentOption {
		return lipgloss.JoinVertical(lipgloss.Left, m.environments.View(), m.help.View(m.keys), m.errors.View())
	}
	if m.focusedView == tagsOption {
		if len(m.tags.Items()) > 0 {
			return lipgloss.JoinVertical(lipgloss.Left, m.tags.View(), m.help.View(m.keys), m.errors.View())
		} else {
			return lipgloss.JoinVertical(lipgloss.Left, "🕦 loading...", m.errors.View())
		}
	}

	deployToStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("205")).Bold(true)
	deployTo := "DEPLOY TO: " + deployToStyle.Render(string(m.selectedEnvironment)) + "\n"

	return lipgloss.JoinVertical(lipgloss.Left, m.CustomTableModel.View(), deployTo, m.help.View(m.keys), m.errors.View())
}

func Deploy() {
	deploymentModel := NewDeploymentModel()
	if _, err := tea.NewProgram(deploymentModel).Run(); err != nil {
		fmt.Printf("Uh oh, there was an error: %v\n", err)
		os.Exit(1)
	}
}
