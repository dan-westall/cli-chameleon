package tui

import (
	"testing"

	"github.com/dan-westall/cli-chameleon/internal/config"
	tea "github.com/charmbracelet/bubbletea"
)

func TestMenuNavigation(t *testing.T) {
	cfg := &config.Config{
		Name: "Test",
		Commands: []config.Command{
			{Name: "build", Description: "Build"},
			{Name: "test", Description: "Test"},
			{Name: "deploy", Description: "Deploy"},
		},
	}

	m := NewModel(cfg)

	if m.cursor != 0 {
		t.Fatal("cursor should start at 0")
	}

	// Move down
	updated, _ := m.Update(tea.KeyMsg{Type: tea.KeyDown})
	m = updated.(Model)
	if m.cursor != 1 {
		t.Fatalf("expected cursor 1, got %d", m.cursor)
	}

	// Move down again
	updated, _ = m.Update(tea.KeyMsg{Type: tea.KeyDown})
	m = updated.(Model)
	if m.cursor != 2 {
		t.Fatalf("expected cursor 2, got %d", m.cursor)
	}

	// Can't go past end
	updated, _ = m.Update(tea.KeyMsg{Type: tea.KeyDown})
	m = updated.(Model)
	if m.cursor != 2 {
		t.Fatalf("expected cursor to stay at 2, got %d", m.cursor)
	}

	// Move up
	updated, _ = m.Update(tea.KeyMsg{Type: tea.KeyUp})
	m = updated.(Model)
	if m.cursor != 1 {
		t.Fatalf("expected cursor 1, got %d", m.cursor)
	}
}

func TestMenuViewRenders(t *testing.T) {
	cfg := &config.Config{
		Name: "My Project",
		Commands: []config.Command{
			{Name: "build", Description: "Build it"},
		},
	}

	m := NewModel(cfg)
	view := m.View()

	if len(view) == 0 {
		t.Fatal("view should not be empty")
	}
}
