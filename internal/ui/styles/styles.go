package styles

import "github.com/charmbracelet/lipgloss"

var MenuStyle = lipgloss.NewStyle().
	Padding(1, 2).
	Border(lipgloss.NormalBorder()).
	Foreground(lipgloss.Color("#FFFFFF"))
var TableStyle = lipgloss.NewStyle().
	Padding(1, 2).
	Border(lipgloss.NormalBorder()).
	Foreground(lipgloss.Color("#FFFFFF"))
var HeaderStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("205")).Bold(true)
var ItemStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("240"))
var OnStyle = ItemStyle.Copy().Foreground(lipgloss.Color("127"))
var ChangedStyle = ItemStyle.Copy().Foreground(lipgloss.Color("9"))
var KeyBindingStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#42A800")).Bold(true)
