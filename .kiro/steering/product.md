# Product: Chameleon

Chameleon is a terminal UI command runner that reads a per-project `chameleon.yaml` and dynamically builds an interactive menu of available commands. It provides a consistent interface across projects regardless of underlying tooling (npm, Python, Make, Docker, etc.).

## Key Capabilities
- Interactive TUI menu built from YAML config
- Single command, sequential, and parallel execution modes
- Live streaming output in a 70/30 split-panel view
- Auto-generates a template config on first run
- Single binary with no runtime dependencies

## User Flow
1. User runs `chameleon` in a project directory
2. If no `chameleon.yaml` exists, a template is created and the program exits
3. If config exists, an interactive menu is displayed
4. User selects a command; it runs either in background (with toast feedback) or in streaming mode
