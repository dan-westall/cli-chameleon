package tui

import (
	"bufio"
	"io"
	"os/exec"
	"strings"

	"github.com/dan-westall/cli-chameleon/internal/config"
	"github.com/dan-westall/cli-chameleon/internal/executor"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type streamState struct {
	cmd     *exec.Cmd
	reader  io.ReadCloser
	output  []string
	name    string
	done    bool
	scroll  int
}

type streamOutputMsg struct {
	line string
}

type streamDoneMsg struct{}

func (m Model) startStream(cmd config.Command) (tea.Model, tea.Cmd) {
	proc, reader, err := executor.StreamCmd(cmd.Run)
	if err != nil {
		m.toast = &toastState{success: false, message: "✗ " + err.Error()}
		m.view = toastView
		return m, nil
	}

	m.stream = &streamState{
		cmd:    proc,
		reader: reader,
		name:   cmd.Name,
	}
	m.view = streamView

	return m, m.readStream(reader)
}

func (m Model) readStream(reader io.ReadCloser) tea.Cmd {
	return func() tea.Msg {
		scanner := bufio.NewScanner(reader)
		if scanner.Scan() {
			return streamOutputMsg{line: scanner.Text()}
		}
		return streamDoneMsg{}
	}
}

func (m Model) updateStream(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "esc", "q", "ctrl+c":
			if m.stream != nil && m.stream.cmd != nil && !m.stream.done {
				m.stream.cmd.Process.Kill()
				m.stream.reader.Close()
			}
			m.stream = nil
			m.view = menuView
			return m, nil
		}
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil
	case streamOutputMsg:
		m.stream.output = append(m.stream.output, msg.line)
		// Auto-scroll to bottom
		visibleLines := m.height - 4
		if len(m.stream.output) > visibleLines {
			m.stream.scroll = len(m.stream.output) - visibleLines
		}
		return m, m.readStream(m.stream.reader)
	case streamDoneMsg:
		m.stream.done = true
		return m, nil
	}
	return m, nil
}

func (m Model) viewStream() string {
	if m.stream == nil {
		return ""
	}

	totalWidth := m.width
	if totalWidth == 0 {
		totalWidth = 80
	}
	totalHeight := m.height
	if totalHeight == 0 {
		totalHeight = 24
	}

	leftWidth := int(float64(totalWidth) * 0.7)
	rightWidth := totalWidth - leftWidth - 1

	// Output panel (70%)
	visibleLines := totalHeight - 4
	start := m.stream.scroll
	end := start + visibleLines
	if end > len(m.stream.output) {
		end = len(m.stream.output)
	}

	var lines []string
	if start < len(m.stream.output) {
		lines = m.stream.output[start:end]
	}

	// Truncate lines to fit width
	for i, l := range lines {
		if len(l) > leftWidth-2 {
			lines[i] = l[:leftWidth-2]
		}
	}

	// Pad to fill height
	for len(lines) < visibleLines {
		lines = append(lines, "")
	}

	outputContent := strings.Join(lines, "\n")
	outputStyle := lipgloss.NewStyle().
		Width(leftWidth).
		Height(visibleLines).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("8"))

	// Info panel (30%)
	status := "● Running"
	statusColor := lipgloss.Color("11")
	if m.stream.done {
		status = "✓ Complete"
		statusColor = lipgloss.Color("10")
	}

	infoContent := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("10")).
		Render(m.stream.name) + "\n\n" +
		lipgloss.NewStyle().Foreground(statusColor).Render(status)

	infoStyle := lipgloss.NewStyle().
		Width(rightWidth).
		Height(visibleLines).
		Padding(1).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("8"))

	// Compose panels
	panels := lipgloss.JoinHorizontal(
		lipgloss.Top,
		outputStyle.Render(outputContent),
		infoStyle.Render(infoContent),
	)

	help := helpStyle.Render("esc back • q quit")

	return panels + "\n" + help
}
