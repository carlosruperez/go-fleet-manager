package menu

import (
	"strings"

	"github.com/go-fleet-manager/internal/ui/styles"
)

func (m MenuModel) View() string {
	var b strings.Builder

	for i, choice := range m.Choices {
		cursor := " "
		style := styles.ItemStyle
		if i == m.Cursor {
			style = styles.OnStyle
			cursor = ">"
		}
		b.WriteString(style.Render(cursor + " " + choice.String()))
		b.WriteString("\n")
	}

	return styles.MenuStyle.Render(b.String())
}
