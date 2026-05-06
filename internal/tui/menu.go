package tui

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	titleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("10")).
			MarginBottom(1)

	itemStyle = lipgloss.NewStyle().
			PaddingLeft(2)

	selectedStyle = lipgloss.NewStyle().
			PaddingLeft(1).
			Foreground(lipgloss.Color("10")).
			Bold(true)

	descStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("8"))

	helpStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("8")).
			MarginTop(1)
)

func (m Model) updateMenu(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}
		case "down", "j":
			if m.cursor < len(m.cfg.Commands)-1 {
				m.cursor++
			}
		case "enter":
			if len(m.cfg.Commands) > 0 {
				return m.executeCommand(m.cfg.Commands[m.cursor])
			}
		}
	}
	return m, nil
}

func (m Model) viewMenu() string {
	var b strings.Builder

	title := "🦎 " + m.cfg.Name
	if m.cfg.Name == "" {
		title = "🦎 Chameleon"
	}
	b.WriteString(titleStyle.Render(title))
	b.WriteString("\n")

	for i, cmd := range m.cfg.Commands {
		name := cmd.Name
		desc := ""
		if cmd.Description != "" {
			desc = descStyle.Render(" — " + cmd.Description)
		}

		if i == m.cursor {
			b.WriteString(selectedStyle.Render("▸ "+name) + desc + "\n")
		} else {
			b.WriteString(itemStyle.Render("  "+name) + desc + "\n")
		}
	}

	if len(m.cfg.Commands) == 0 {
		b.WriteString(descStyle.Render("  No commands configured. Edit chameleon.yaml to add commands.\n"))
	}

	help := fmt.Sprintf("\n%s", helpStyle.Render("↑/↓ navigate • enter select • q quit"))
	b.WriteString(help)

	return b.String()
}
