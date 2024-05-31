package ui

import (
	"context"
	"fmt"
	"os"
	"reflect"
	"time"

	"github.com/go-redis/redis/v8"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/go-fleet-manager/config"
	uiCommon "github.com/go-fleet-manager/internal/ui/common"
	uiCustomTable "github.com/go-fleet-manager/internal/ui/customTable"
	uiErrorsModel "github.com/go-fleet-manager/internal/ui/errors"
)

type RemoveCacheModel struct {
	uiCustomTable.CustomTableModel
	keys   removeCacheKeyMap
	help   help.Model
	errors uiErrorsModel.ErrorsModel
	items  []removeCacheItem
}

type removeCacheItem struct {
	name        string
	cacheConfig config.CacheConfig
	results     string
}

type removeCacheKeyMap struct {
	uiCommon.KeyMap
}

func (k removeCacheKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Help, k.Esc}
}

func (k removeCacheKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Up, k.Down},
		{k.Enter, k.Help},
		{k.Esc, k.Quit},
	}
}

var removeCacheKeys = removeCacheKeyMap{
	KeyMap: uiCommon.Keys,
}

func (model *RemoveCacheModel) Equals(other *RemoveCacheModel) bool {
	return reflect.DeepEqual(model.items, other.items)
}

func NewRemoveCacheModel() RemoveCacheModel {

	caches := config.GetCaches()

	items := []removeCacheItem{}

	for _, cache := range caches {
		items = append(items, removeCacheItem{name: cache.Name, cacheConfig: cache, results: ""})
	}

	columns := []table.Column{
		{Title: "HOST", Width: 25},
		{Title: "results", Width: 15},
	}

	values := []table.Row{}
	for _, item := range items {
		values = append(values, table.Row{item.name})
	}

	customTableModel := uiCustomTable.NewCustomTableModel(columns, values)

	m := RemoveCacheModel{
		keys:             removeCacheKeys,
		help:             help.New(),
		errors:           uiErrorsModel.NewErrorsModel(),
		items:            items,
		CustomTableModel: customTableModel,
	}
	return m
}

type flushMsg struct {
	index   int
	results string
	err     error
}

func flushallRedis(index int, cacheConfig config.CacheConfig) tea.Cmd {
	return func() tea.Msg {
		rdb := redis.NewClient(&redis.Options{
			Addr: cacheConfig.Host + ":" + cacheConfig.Port,
		})

		ctx := context.Background()
		ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
		defer cancel()

		err := rdb.FlushAll(ctx).Err()
		if err != nil {
			return flushMsg{index: index, results: "👎", err: err}
		}

		return flushMsg{index: index, results: "✅"}
	}
}

func (m *RemoveCacheModel) updateTable() {
	values := []table.Row{}
	for _, item := range m.items {
		values = append(values, table.Row{item.name, item.results})
	}
	m.SetRows(values)
}

func (m RemoveCacheModel) Init() tea.Cmd {
	return nil
}

func (m RemoveCacheModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case uiCommon.TickMsg:
		for i := range m.items {
			m.items[i].results = ""
		}
		m.updateTable()
		return m, nil
	case flushMsg:
		m.items[msg.index].results = msg.results
		if msg.err != nil {
			m.items[msg.index].results = "👎"
			m.errors.AddError(msg.err)
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
			return m, tea.Batch(flushallRedis(m.Cursor(), m.items[m.Cursor()].cacheConfig), uiCommon.UpdateEverySeconds(10))
		case key.Matches(msg, m.keys.Help):
			m.help.ShowAll = !m.help.ShowAll
		}
	}
	return m, nil
}

func (m RemoveCacheModel) View() string {
	return lipgloss.JoinVertical(lipgloss.Left, m.CustomTableModel.View(), m.help.View(m.keys), m.errors.View())
}

func RemoveAllCache() {

	removeCacheModel := NewRemoveCacheModel()

	if _, err := tea.NewProgram(removeCacheModel).Run(); err != nil {
		fmt.Printf("Uh oh, there was an error: %v\n", err)
		os.Exit(1)
	}
}
