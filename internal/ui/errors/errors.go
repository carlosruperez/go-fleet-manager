package errors

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type ErrorsModel struct {
	errors []error
}

func NewErrorsModel() ErrorsModel {
	return ErrorsModel{}
}

func (m *ErrorsModel) AddError(err error) {
	m.errors = append(m.errors, err)
}

func (m ErrorsModel) Init() tea.Cmd {
	return nil
}

func (m ErrorsModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return m, nil
}

var baseStyle = lipgloss.NewStyle().
	BorderStyle(lipgloss.NormalBorder()).
	BorderForeground(lipgloss.Color("#D46A6A"))

func (m ErrorsModel) View() string {
	if len(m.errors) > 0 {
		errorsString := "ERRORS:"
		lastNErrors := 5
		if len(m.errors)-1 < 5 {
			lastNErrors = len(m.errors)
		}
		for _, element := range m.errors[len(m.errors)-lastNErrors:] {
			errorsString += "\n" + element.Error()
		}
		return baseStyle.Render(errorsString)
	}
	return ""
}
