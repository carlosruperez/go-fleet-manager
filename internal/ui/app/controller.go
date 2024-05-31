package app

import (
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
	"os"
)

func RunApp() {
	m := newAppModel()
	if _, err := tea.NewProgram(m).Run(); err != nil {
		fmt.Printf("Uh oh, there was an error: %v\n", err)
		os.Exit(1)
	}
}
