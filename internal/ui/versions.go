package ui

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"reflect"
	"sort"
	"strings"
	"time"

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
	"github.com/go-fleet-manager/internal/usecases"
)

var secondsBetweenUpdates = 10

type VersionsModel struct {
	uiCustomTable.CustomTableModel
	keys   versionsKeyMap
	help   help.Model
	errors uiErrorsModel.ErrorsModel
	items  []item // TODO use Rows() and SetRows() ...
}

type item struct {
	name    string
	urls    []string
	results []string
}

type versionsKeyMap struct {
	uiCommon.KeyMap
}

func (k versionsKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Up, k.Down, k.Esc, k.Quit}
}

var versionsKeys = versionsKeyMap{
	KeyMap: uiCommon.Keys,
}

func (model *VersionsModel) Equals(other *VersionsModel) bool {
	return reflect.DeepEqual(model.items, other.items)
}

func NewVersionsModel() VersionsModel {

	urls := []string{
		"http://carlosruperez.internal-api.dev.internal/",
		"http://carlosruperez.internal-api.pre.internal/",
		"http://carlosruperez.internal-api.internal/",
	}

	items := []item{}
	servicesPaths := getServicesPaths()
	for _, servicePath := range servicesPaths {
		ms_urls := []string{}
		for _, url := range urls {
			ms_urls = append(ms_urls, url+servicePath+"/version/")
		}
		items = append(items, item{name: servicePath, urls: ms_urls, results: []string{"", "", ""}})
	}

	columns := []table.Column{
		{Title: "n°", Width: 5},
		{Title: "URL", Width: 25},
		{Title: "DEV", Width: 10},
		{Title: "PRE", Width: 10},
		{Title: "PROD", Width: 18},
	}

	values := []table.Row{}
	for i, item := range items {
		values = append(values, table.Row{fmt.Sprint(i + 1), item.name, item.results[0], item.results[1], item.results[2]})
	}
	customTableModel := uiCustomTable.NewCustomTableModel(columns, values)

	m := VersionsModel{
		keys:             versionsKeys,
		help:             help.New(),
		errors:           uiErrorsModel.NewErrorsModel(),
		items:            items,
		CustomTableModel: customTableModel}

	return m
}

type reqContent struct {
	Version string `json:"version"`
}

type fetchMsg struct {
	index       int
	environment int
	statusCode  int
	err         error
	version     string
}

func fetch(index int, environment int, url string) tea.Cmd {
	return func() tea.Msg {
		client := http.Client{
			Timeout: 5 * time.Second,
		}
		resp, err := client.Get(url)
		if err != nil {
			return fetchMsg{index: index, environment: environment, err: err}
		}
		defer resp.Body.Close()

		content := reqContent{}
		err = json.NewDecoder(resp.Body).Decode(&content)

		return fetchMsg{index: index, environment: environment, statusCode: resp.StatusCode, version: content.Version, err: err}
	}
}

func (m *VersionsModel) updateTable() {
	values := []table.Row{}
	for i, item := range m.items {
		values = append(values, table.Row{fmt.Sprint(i + 1), item.name, item.results[0], item.results[1], item.results[2]})
	}
	m.SetRows(values)
}

func (m VersionsModel) Init() tea.Cmd {
	cmds := make([]tea.Cmd, 0)
	for i := range m.items {
		for j, url := range m.items[i].urls {
			cmds = append(cmds, fetch(i, j, url))
		}
	}
	cmds = append(cmds, uiCommon.UpdateEverySeconds(secondsBetweenUpdates))
	return tea.Batch(cmds...)
}

func (m VersionsModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {

	switch msg := msg.(type) {
	case fetchMsg:
		m.items[msg.index].results[msg.environment] = msg.version
		if msg.err != nil {
			m.items[msg.index].results[msg.environment] = "👎"
			m.errors.AddError(msg.err)
		}
		if m.items[msg.index].results[1] != "" && m.items[msg.index].results[1] != "👎" && m.items[msg.index].results[1] != "👀" {
			if m.items[msg.index].results[1] != m.items[msg.index].results[2] {
				var sb strings.Builder
				sb.WriteString("\033[31m")
				sb.WriteString(m.items[msg.index].results[2])
				sb.WriteString("\033[0m")
				m.items[msg.index].results[2] = sb.String()
			}
		}

		m.updateTable()
		return m, nil
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
			return m, nil
		case key.Matches(msg, m.keys.Esc):
			return m, nil
		}
	case uiCommon.TickMsg:
		cmds := make([]tea.Cmd, 0)
		for i := range m.items {
			m.items[i].results = []string{"👀", "👀", "👀"}
			for j, url := range m.items[i].urls {
				cmds = append(cmds, fetch(i, j, url))
			}
		}
		cmds = append(cmds, uiCommon.UpdateEverySeconds(secondsBetweenUpdates))
		return m, tea.Batch(cmds...)
	}
	return m, nil
}

func (m VersionsModel) View() string {
	return lipgloss.JoinVertical(lipgloss.Left, m.CustomTableModel.View(), m.help.View(m.keys), m.errors.View())
}

func getServicesPaths() []string {
	repositories := usecases.GetRepositories()
	msRepositories := []repository.Repository{}
	for _, repository := range repositories {
		if repository.Type == common.Microservice {
			msRepositories = append(msRepositories, repository)
		}
	}

	servicesPaths := []string{}
	for _, msRepository := range msRepositories {
		msPath, err := msRepository.GetMSPath()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		servicesPaths = append(servicesPaths, msPath)
	}
	sort.Strings(servicesPaths)

	return servicesPaths
}

func GetVersions() {
	versionsModel := NewVersionsModel()
	if _, err := tea.NewProgram(versionsModel).Run(); err != nil {
		fmt.Printf("Uh oh, there was an error: %v\n", err)
		os.Exit(1)
	}
}
