# Project Structure

```
chameleon/
├── main.go                          # Entry point, version flag, config bootstrap
├── Makefile                         # Build targets (build, install, test, clean)
├── go.mod / go.sum                  # Go module definition
├── bin/                             # Build output (gitignored)
└── internal/                        # Private packages (not importable externally)
    ├── config/
    │   ├── config.go                # Config struct, YAML parsing, Load/Exists
    │   ├── config_test.go           # Config loading tests
    │   ├── template.go              # Template YAML generation for first run
    │   └── template_test.go         # Template tests
    ├── executor/
    │   ├── executor.go              # Command execution: sequential, parallel, streaming
    │   └── executor_test.go         # Executor tests
    └── tui/
        ├── model.go                 # Root Bubble Tea model, view routing
        ├── model_test.go            # Model/navigation tests
        ├── menu.go                  # Main menu view and key handling
        ├── stream.go                # 70/30 split streaming view
        └── toast.go                 # Success/failure toast overlay
```

## Architecture Patterns
- **Elm architecture** via Bubble Tea: Model → Update → View cycle
- **View states**: `menuView`, `toastView`, `streamView` — routed in the root model
- **Internal packages only**: all code lives under `internal/` to prevent external imports
- **Separation of concerns**: config parsing, command execution, and TUI are independent packages
- **Custom YAML unmarshalling**: `Command.UnmarshalYAML` handles polymorphic `run` field

## Conventions
- No exported API — this is a standalone CLI tool
- Tests live alongside source files (`_test.go` suffix)
- Styles defined as package-level `lipgloss` variables in the view files
- Commands executed via `sh -c` for shell expansion support
