# Tech Stack

## Language & Runtime
- Go 1.25+
- Module path: `github.com/bsport/chameleon`

## Key Libraries
- **Bubble Tea** (`github.com/charmbracelet/bubbletea`) — TUI framework (Elm architecture)
- **Lip Gloss** (`github.com/charmbracelet/lipgloss`) — TUI styling
- **Bubbles** (`github.com/charmbracelet/bubbles`) — TUI components
- **gopkg.in/yaml.v3** — YAML parsing for config

## Build System
- **Makefile** with the following targets:
  - `make build` — compile binary to `bin/chameleon`
  - `make install` — install to `$GOPATH/bin`
  - `make test` — run all tests with `go test ./... -v`
  - `make clean` — remove build artifacts
- Version injected via `-ldflags` from git tags

## Testing
- Standard Go testing (`testing` package)
- No external test framework
- Run tests: `go test ./... -v`

## Configuration Format
- YAML (`chameleon.yaml`) in the working directory
- Custom `UnmarshalYAML` for flexible `run` field (string or string array)
