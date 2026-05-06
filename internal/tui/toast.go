package tui

import (
	"time"

	"github.com/dan-westall/cli-chameleon/internal/config"
	"github.com/dan-westall/cli-chameleon/internal/executor"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type toastState struct {
	success bool
	message string
}

type cmdDoneMsg struct {
	result executor.Result
}

type toastDismissMsg struct{}

func (m Model) executeCommand(cmd config.Command) (tea.Model, tea.Cmd) {
	if cmd.Stream {
		return m.startStream(cmd)
	}

	// Run synchronously via tea.Cmd
	return m, func() tea.Msg {
		result := executor.Run(cmd.Run, cmd.Parallel)
		return cmdDoneMsg{result: result}
	}
}

func (m Model) updateToast(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg.(type) {
	case tea.KeyMsg:
		m.view = menuView
		m.toast = nil
		return m, nil
	case toastDismissMsg:
		m.view = menuView
		m.toast = nil
		return m, nil
	}
	return m, nil
}

func (m Model) handleCmdDone(result executor.Result) (tea.Model, tea.Cmd) {
	success := result.Err == nil
	msg := "✓ Success"
	if !success {
		msg = "✗ Failed: " + result.Err.Error()
	}
	m.toast = &toastState{success: success, message: msg}
	m.view = toastView

	// Auto-dismiss after 2 seconds
	return m, tea.Tick(2*time.Second, func(time.Time) tea.Msg {
		return toastDismissMsg{}
	})
}

func (m Model) viewToast() string {
	style := lipgloss.NewStyle().
		Bold(true).
		Padding(1, 2).
		MarginTop(1)

	if m.toast.success {
		style = style.Foreground(lipgloss.Color("10"))
	} else {
		style = style.Foreground(lipgloss.Color("9"))
	}

	return m.viewMenu() + "\n\n" + style.Render(m.toast.message)
}
