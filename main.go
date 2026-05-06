package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/dan-westall/cli-chameleon/internal/config"
	"github.com/dan-westall/cli-chameleon/internal/tui"
)

var version = "dev"

func main() {
	if len(os.Args) > 1 && os.Args[1] == "--version" {
		fmt.Printf("chameleon %s\n", version)
		os.Exit(0)
	}

	if !config.Exists() {
		if err := config.CreateTemplate(); err != nil {
			fmt.Fprintf(os.Stderr, "Error creating %s: %v\n", config.FileName, err)
			os.Exit(1)
		}
		fmt.Printf("Created %s — edit it to configure your commands.\n", config.FileName)
		os.Exit(0)
	}

	cfg, err := config.Load()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	p := tea.NewProgram(tui.NewModel(cfg), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
