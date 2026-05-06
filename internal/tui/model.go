package tui

import (
	"github.com/dan-westall/cli-chameleon/internal/config"
	tea "github.com/charmbracelet/bubbletea"
)

type view int

const (
	menuView view = iota
	toastView
	streamView
)

type Model struct {
	cfg      *config.Config
	cursor   int
	view     view
	width    int
	height   int
	toast    *toastState
	stream   *streamState
}

func NewModel(cfg *config.Config) Model {
	return Model{cfg: cfg}
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		if m.stream != nil {
			return m.updateStream(msg)
		}
		return m, nil
	case cmdDoneMsg:
		return m.handleCmdDone(msg.result)
	case streamOutputMsg:
		return m.updateStream(msg)
	case streamDoneMsg:
		return m.updateStream(msg)
	}

	switch m.view {
	case menuView:
		return m.updateMenu(msg)
	case toastView:
		return m.updateToast(msg)
	case streamView:
		return m.updateStream(msg)
	}
	return m, nil
}

func (m Model) View() string {
	switch m.view {
	case toastView:
		return m.viewToast()
	case streamView:
		return m.viewStream()
	default:
		return m.viewMenu()
	}
}
