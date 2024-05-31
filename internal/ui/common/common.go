package common

import (
	tea "github.com/charmbracelet/bubbletea"
	"time"
)

type TickMsg struct{}

func UpdateEverySeconds(seconds int) tea.Cmd {
	return tea.Every(time.Duration(seconds)*time.Second, func(t time.Time) tea.Msg {
		return TickMsg{}
	})
}
