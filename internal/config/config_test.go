package config

import (
	"os"
	"testing"
)

func TestLoadStringRun(t *testing.T) {
	writeTestConfig(t, `
name: "Test Project"
commands:
  - name: build
    description: "Build it"
    run: "go build ./..."
`)
	cfg, err := Load()
	if err != nil {
		t.Fatal(err)
	}
	if cfg.Name != "Test Project" {
		t.Fatalf("expected name 'Test Project', got %q", cfg.Name)
	}
	if len(cfg.Commands) != 1 {
		t.Fatalf("expected 1 command, got %d", len(cfg.Commands))
	}
	cmd := cfg.Commands[0]
	if len(cmd.Run) != 1 || cmd.Run[0] != "go build ./..." {
		t.Fatalf("unexpected run: %v", cmd.Run)
	}
	if cmd.Stream || cmd.Parallel {
		t.Fatal("expected stream and parallel to be false")
	}
}

func TestLoadArrayRun(t *testing.T) {
	writeTestConfig(t, `
name: "Multi"
commands:
  - name: deploy
    description: "Deploy"
    run:
      - "docker build ."
      - "docker push app"
    parallel: true
`)
	cfg, err := Load()
	if err != nil {
		t.Fatal(err)
	}
	cmd := cfg.Commands[0]
	if len(cmd.Run) != 2 {
		t.Fatalf("expected 2 run items, got %d", len(cmd.Run))
	}
	if !cmd.Parallel {
		t.Fatal("expected parallel to be true")
	}
}

func TestLoadStreamFlag(t *testing.T) {
	writeTestConfig(t, `
name: "Stream"
commands:
  - name: ingest
    description: "Ingest data"
    run: "python ingest.py"
    stream: true
`)
	cfg, err := Load()
	if err != nil {
		t.Fatal(err)
	}
	if !cfg.Commands[0].Stream {
		t.Fatal("expected stream to be true")
	}
}

func TestLoadInvalidYAML(t *testing.T) {
	writeTestConfig(t, `invalid: [`)
	_, err := Load()
	if err == nil {
		t.Fatal("expected error for invalid YAML")
	}
}

func writeTestConfig(t *testing.T, content string) {
	t.Helper()
	if err := os.WriteFile(FileName, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { os.Remove(FileName) })
}
