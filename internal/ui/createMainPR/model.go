package createMainPR

import (
	"context"
	"fmt"
	"reflect"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
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
	branchesOption
)

type CreateMainPRModel struct {
	uiCustomTable.CustomTableModel
	keys         createMainPRKeyMap
	help         help.Model
	errors       uiErrorsModel.ErrorsModel
	items        []createMainPRItem
	selectedItem int
	branches     list.Model
	focusedView  viewOption
}

type branchItem struct {
	name string
	list.Item
}

func (b branchItem) Title() string       { return b.name }
func (b branchItem) Description() string { return b.name }
func (b branchItem) FilterValue() string { return b.name }

type createMainPRKeyMap struct {
	uiCommon.KeyMap
	Run    key.Binding
	Select key.Binding
}

func (k createMainPRKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Up, k.Down, k.Select, k.Run, k.Enter, k.Esc, k.Quit}
}

var createMainPRKeys = createMainPRKeyMap{
	KeyMap: uiCommon.Keys,
	Run: key.NewBinding(
		key.WithKeys("r"),
		key.WithHelp(styles.KeyBindingStyle.Render("r"), "Run selected option"),
	),
	Select: key.NewBinding(
		key.WithKeys("s"),
		key.WithHelp(styles.KeyBindingStyle.Render("s"), "Select branch"),
	),
}

type createMainPRItem struct {
	repository     repository.Repository
	selectedBranch string
	result         string
}

func (model *CreateMainPRModel) Equals(other *CreateMainPRModel) bool {
	return reflect.DeepEqual(model.items, other.items)
}

func (m *CreateMainPRModel) updateTable() {
	values := []table.Row{}
	for i, item := range m.items {
		values = append(values, table.Row{fmt.Sprint(i + 1), item.repository.Name, item.selectedBranch, item.result})
	}
	m.SetRows(values)
}

func (m *CreateMainPRModel) resetResults() {
	for i := range m.items {
		m.items[i].result = ""
	}
	m.updateTable()
}

type createMainPRMsg struct {
	index  int
	repo   repository.Repository
	result string
	err    error
}

func createMainPR(index int, repo repository.Repository, branch string) tea.Cmd {
	ctx, cancel := context.WithCancel(context.Background()) // TODO: use context to get feedback
	go func() tea.Msg {
		<-ctx.Done()
		// The context was cancelled, stop the update process
		// and return a special message
		return createMainPRMsg{index, repo, "❌" + "Cancelled", nil}
	}()
	return func() tea.Msg {
		_, _, err := usecases.CreateMainPR(repo, branch, ctx)
		if err != nil {
			cancel()
			return createMainPRMsg{index, repo, "❌", err}
		}
		return createMainPRMsg{index, repo, "🚀", nil}
	}
}

type selectRepositoryBranchMsg struct {
	selectedItem int
	branches     []string
	err          error
}

func initSelectRepositoryBranch(repo repository.Repository, selectedItem int) tea.Cmd {
	ctx, cancel := context.WithCancel(context.Background()) // TODO: use context to get feedback
	go func() tea.Msg {
		<-ctx.Done()
		// The context was cancelled, stop the update process
		// and return a special message
		return selectRepositoryBranchMsg{selectedItem, []string{"❌" + "Cancelled"}, nil}
	}()
	return func() tea.Msg {
		branches, err := usecases.GetRepositoryBranches(repo, ctx)
		if err != nil {
			cancel()
			return selectRepositoryBranchMsg{selectedItem, nil, err}
		}
		return selectRepositoryBranchMsg{selectedItem, branches, nil}
	}
}

func NewCreateMainPRModel() CreateMainPRModel {

	items := []createMainPRItem{}

	repositories := usecases.GetRepositories()

	for _, repo := range repositories {
		if repo.Type != common.Microservice {
			continue
		}
		items = append(items, createMainPRItem{repository: repo}) // TODO hardcoded
	}

	columns := []table.Column{
		{Title: "n°", Width: 5},
		{Title: "Repository", Width: 50},
		{Title: "Selected Branch", Width: 25},
		{Title: "Result", Width: 15},
	}

	values := []table.Row{}
	for i, item := range items {
		values = append(values, table.Row{fmt.Sprint(i + 1), item.repository.Name, item.result})
	}
	customTableModel := uiCustomTable.NewCustomTableModel(columns, values)

	delegate := list.NewDefaultDelegate()
	delegate.ShowDescription = false
	newList := list.New([]list.Item{}, delegate, 50, 13)
	newList.Title = "Select a branch to create the PR from:"
	newList.SetShowHelp(false)
	newList.SetShowFilter(false)
	newList.SetShowPagination(false)
	newList.SetShowStatusBar(false)

	m := CreateMainPRModel{
		keys:             createMainPRKeys,
		help:             help.New(),
		errors:           uiErrorsModel.NewErrorsModel(),
		items:            items,
		CustomTableModel: customTableModel,
		branches:         newList,
		focusedView:      tableOption,
	}

	return m
}
