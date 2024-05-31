package createMainPR

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
)

func CreateMainPR() {
	createMainPRModel := NewCreateMainPRModel()
	if _, err := tea.NewProgram(createMainPRModel).Run(); err != nil {
		fmt.Printf("Uh oh, there was an error: %v\n", err)
		os.Exit(1)
	}
}
