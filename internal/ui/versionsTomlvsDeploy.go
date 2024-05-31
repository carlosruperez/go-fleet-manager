package ui

import (
	"encoding/base64"
	"encoding/json"
	"strings"

	"fmt"
	"net/http"
	"os"

	"reflect"
	"sort"
	"time"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/cli/go-gh/v2"
	"github.com/pelletier/go-toml"

	"github.com/go-fleet-manager/config"
	"github.com/go-fleet-manager/internal/common"
	"github.com/go-fleet-manager/internal/repository"
	uiCommon "github.com/go-fleet-manager/internal/ui/common"
	uiCustomTable "github.com/go-fleet-manager/internal/ui/customTable"
	"github.com/go-fleet-manager/internal/usecases"
)

type VersionsTomlVsDeployModel struct {
	uiCustomTable.CustomTableModel
	keys                versionsTomlVsDeployKeyMap
	help                help.Model
	enter               bool
	items               []versionsTomlVsDeployItem
	versionsDeployModel versionsDeployModel
}

type versionsTomlVsDeployItem struct {
	name string
	url  string
}

type versionsTomlVsDeployKeyMap struct {
	uiCommon.KeyMap
}

func (k versionsTomlVsDeployKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Up, k.Down, k.Enter, k.Esc, k.Quit}
}

func (k versionsTomlVsDeployKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Up, k.Down},
		{k.Enter, k.Help},
		{k.Esc, k.Quit},
	}
}

var versionsTomlVsDeployKeys = versionsTomlVsDeployKeyMap{
	KeyMap: uiCommon.Keys,
}

func (model *VersionsTomlVsDeployModel) Equals(other *VersionsTomlVsDeployModel) bool {
	return reflect.DeepEqual(model.items, other.items)
}

func NewVersionsTomlVsDeployModel() VersionsTomlVsDeployModel {

	environments := config.GetEnvironments()

	items := []versionsTomlVsDeployItem{}

	for _, environment := range environments {
		items = append(items, versionsTomlVsDeployItem{name: environment.Name, url: environment.Url})
	}

	columns := []table.Column{
		{Title: "Environment", Width: 25},
	}

	values := []table.Row{}
	for _, item := range items {
		values = append(values, table.Row{item.name})
	}

	customTableModel := uiCustomTable.NewCustomTableModel(columns, values)

	m := VersionsTomlVsDeployModel{
		keys:             versionsTomlVsDeployKeys,
		help:             help.New(),
		items:            items,
		CustomTableModel: customTableModel,
	}
	return m
}

func (m VersionsTomlVsDeployModel) Init() tea.Cmd {
	return nil
}

func (m VersionsTomlVsDeployModel) forwardMsgToViewDeploy(msg tea.Msg) (tea.Model, tea.Cmd) {
	if m.enter {
		if m.versionsDeployModel.Equals(&versionsDeployModel{}) {
			m.versionsDeployModel = newDeployVersionsModel(m.items[m.Cursor()].url)
		}
		newModel, cmd := m.versionsDeployModel.Update(msg)
		m.versionsDeployModel = newModel.(versionsDeployModel)

		return m, cmd
	}
	return m, nil

}
func (m VersionsTomlVsDeployModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.keys.Quit):
			return m, tea.Quit

		case key.Matches(msg, m.keys.Up):
			if m.Cursor() > 0 {
				m.SetCursor(m.Cursor() - 1)
			}
			return m.forwardMsgToViewDeploy(msg)

		case key.Matches(msg, m.keys.Down):
			if m.Cursor() < len(m.items)-1 {
				m.SetCursor(m.Cursor() + 1)
			}
			return m.forwardMsgToViewDeploy(msg)

		case key.Matches(msg, m.keys.Enter):

			if !m.enter {
				m.enter = true
				model, _ := m.forwardMsgToViewDeploy(msg)
				return model, m.versionsDeployModel.Init()
			}
			return m.forwardMsgToViewDeploy(msg)

		case key.Matches(msg, m.keys.Esc):

			return m, nil

		case key.Matches(msg, m.keys.Help):
			m.help.ShowAll = !m.help.ShowAll
			return m.forwardMsgToViewDeploy(msg)

		default:
			return m.forwardMsgToViewDeploy(msg)
		}
	}
	return m.forwardMsgToViewDeploy(msg)
}

func (m VersionsTomlVsDeployModel) View() string {
	if m.enter {
		return lipgloss.JoinHorizontal(lipgloss.Left, m.CustomTableModel.View(), m.versionsDeployModel.View())
	}
	return m.CustomTableModel.View() + m.help.View(m.keys)
}

func CheckVersions() {

	versionsTomlVsDeployModel := NewVersionsTomlVsDeployModel()

	if _, err := tea.NewProgram(versionsTomlVsDeployModel).Run(); err != nil {
		fmt.Printf("Uh oh, there was an error: %v\n", err)
		os.Exit(1)
	}
}

var secondsUpdates = 10

type versionsDeployModel struct {
	uiCustomTable.CustomTableModel
	keys  versionsKeyDeployMap
	help  help.Model
	items []itemDepoly
}

type itemDepoly struct {
	name       string
	deploy_url string
	toml_url   string
	results    []string
}

type versionsKeyDeployMap struct {
	uiCommon.KeyMap
}

func (k versionsKeyDeployMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Up, k.Down, k.Enter, k.Esc, k.Quit}
}

var versionsDeployKeys = versionsKeyDeployMap{
	KeyMap: uiCommon.Keys,
}

func (model *versionsDeployModel) Equals(other *versionsDeployModel) bool {
	return reflect.DeepEqual(model.items, other.items)
}

func newDeployVersionsModel(url string) versionsDeployModel {

	items := []itemDepoly{}
	paths := getDeployServicesPaths()

	for _, path := range paths {

		var ms_deploy_url string
		var ms_toml_url string
		var branch string

		if strings.Contains(url, ".dev.") {
			branch = "develop"
		} else if strings.Contains(url, ".pre.") {
			branch = "release-1.5.x"
		} else {
			branch = "main"
		}

		ms_deploy_url = url + path.ServicePath + "/version"
		ms_toml_url = "/repos/carlosruperez/" + path.MsPath + "/contents/pyproject.toml?ref=" + branch

		items = append(items, itemDepoly{name: path.ServicePath, deploy_url: ms_deploy_url, toml_url: ms_toml_url, results: []string{"", ""}})
	}

	columns := []table.Column{
		{Title: "URL", Width: 25},
		{Title: "TOML", Width: 10},
		{Title: "DEPLOYED", Width: 10},
	}

	values := []table.Row{}
	for _, item := range items {
		values = append(values, table.Row{item.name, item.results[0], item.results[1]})
	}
	customTableModel := uiCustomTable.NewCustomTableModel(columns, values)

	m := versionsDeployModel{
		keys:             versionsDeployKeys,
		help:             help.New(),
		items:            items,
		CustomTableModel: customTableModel}

	return m
}

type reqDeployContent struct {
	Version string `json:"version"`
}

type fetchDeployMsg struct {
	index       int
	environment int
	statusCode  int
	err         error
	version     string
}

type contentDeployData struct {
	Content string `json:"content"`
}

func fetchToml(index int, environment int, url string) tea.Cmd {
	return func() tea.Msg {
		outputBytes, _, err := gh.Exec("api", url)
		if err != nil {
			return fetchDeployMsg{index: index, environment: environment, statusCode: 404, version: ""}
		}

		var contents contentDeployData
		if err := json.Unmarshal(outputBytes.Bytes(), &contents); err != nil {
			return fetchMsg{index: index, environment: environment, err: err}
		}

		decodedContent, err := base64.StdEncoding.DecodeString(contents.Content)
		if err != nil {
			return fetchMsg{index: index, environment: environment, err: err}
		}

		var config map[string]interface{}
		if err := toml.Unmarshal([]byte(decodedContent), &config); err != nil {
			return fetchMsg{index: index, environment: environment, err: err}
		}

		devDependencies := config["tool"].(map[string]interface{})["poetry"].(map[string]interface{})
		version := devDependencies["version"].(string)
		return fetchDeployMsg{index: index, environment: environment, statusCode: 200, version: version}
	}
}

func fetchDeploy(index int, environment int, url string) tea.Cmd {
	return func() tea.Msg {
		client := http.Client{
			Timeout: 5 * time.Second,
		}
		resp, err := client.Get(url)
		if err != nil {
			return fetchDeployMsg{index: index, environment: environment, err: err}
		}
		defer resp.Body.Close()

		content := reqDeployContent{}
		err = json.NewDecoder(resp.Body).Decode(&content)

		if err != nil {
			return fetchDeployMsg{index: index, environment: environment, err: err}
		}

		return fetchDeployMsg{index: index, environment: environment, statusCode: resp.StatusCode, version: content.Version}

	}
}

func (m *versionsDeployModel) updateTable() {
	values := []table.Row{}
	for _, item := range m.items {
		values = append(values, table.Row{item.name, item.results[0], item.results[1]})
	}
	m.SetRows(values)
}

func (m versionsDeployModel) Init() tea.Cmd {
	cmds := make([]tea.Cmd, 0)

	for i := range m.items {
		cmds = append(cmds, fetchToml(i, 0, m.items[i].toml_url))
		cmds = append(cmds, fetchDeploy(i, 1, m.items[i].deploy_url))
	}
	cmds = append(cmds, uiCommon.UpdateEverySeconds(secondsBetweenUpdates))
	return tea.Batch(cmds...)
}

func (m versionsDeployModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case fetchDeployMsg:
		m.items[msg.index].results[msg.environment] = msg.version
		if msg.err != nil {
			m.items[msg.index].results[msg.environment] = "👎"
		}

		if m.items[msg.index].results[0] != "" && m.items[msg.index].results[0] != "👎" && m.items[msg.index].results[0] != "👀" {
			if m.items[msg.index].results[0] != m.items[msg.index].results[1] {
				var sb strings.Builder
				sb.WriteString("\033[31m")
				sb.WriteString(m.items[msg.index].results[1])
				sb.WriteString("\033[0m")
				m.items[msg.index].results[1] = sb.String()
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
			m.items[i].results = []string{"👀", "👀"}
			cmds = append(cmds, fetchToml(i, 0, m.items[i].toml_url))
			cmds = append(cmds, fetchDeploy(i, 1, m.items[i].deploy_url))
		}
		cmds = append(cmds, uiCommon.UpdateEverySeconds(secondsBetweenUpdates))
		return m, tea.Batch(cmds...)
	}
	return m, nil
}

func (m versionsDeployModel) View() string {
	return m.CustomTableModel.View() + m.help.View(m.keys)
}

type PathPair struct {
	ServicePath string
	MsPath      string
}

func getDeployServicesPaths() []PathPair {
	repositories := usecases.GetRepositories()
	msRepositories := []repository.Repository{}
	for _, repository := range repositories {
		if repository.Type == common.Microservice {
			msRepositories = append(msRepositories, repository)
		}
	}

	pathPairs := []PathPair{}

	for _, msRepository := range msRepositories {
		msPath, err := msRepository.GetMSPath()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		pathPairs = append(pathPairs, PathPair{ServicePath: msPath, MsPath: msRepository.Name})
	}
	sort.Slice(pathPairs, func(i, j int) bool {
		return pathPairs[i].ServicePath < pathPairs[j].ServicePath
	})

	return pathPairs
}
