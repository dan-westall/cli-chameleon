# 🦎 Chameleon

One CLI. Every project. Chameleon adapts.

Chameleon is a terminal UI command runner that reads a per-project `chameleon.yaml` and dynamically builds an interactive menu of your commands. Whether a project uses npm, Python, Bash, Make, or Docker — you only need to remember one interface.

```
┌─────────────────────────────────────┐
│ 🦎 My Project                       │
│                                     │
│ ▸ ingest  — Run ingestion pipeline  │
│   build   — Build the project       │
│   deploy  — Deploy to production    │
│   test    — Run test suite          │
│                                     │
│ ↑/↓ navigate • enter select • q quit│
└─────────────────────────────────────┘
```

## What Problem Does This Solve?

Developers jump between projects constantly — and in the AI era, that's never been more true. We're all building prototypes faster than ever, spinning up new tools, experimenting with new stacks. The days of living in one or two languages are gone. With AI-assisted development you might have a Go service, a Python data pipeline, a Node.js frontend, and a Rust CLI all in active rotation. Each project has its own tooling, its own scripts, its own way of doing things. One uses `npm run build`, another uses `make build`, another uses `python scripts/build.py`. The cognitive overhead of remembering project-specific commands across all of these adds up fast.

Chameleon provides a consistent abstraction layer. You define a `chameleon.yaml` in each project that maps familiar command names to project-specific scripts. The TUI dynamically renders whatever commands are available, so you never need to remember which tool a particular project uses.

## Tech Stack

- **Go** — single binary, no runtime dependencies
- **[Bubble Tea](https://github.com/charmbracelet/bubbletea)** — terminal UI framework
- **[Lip Gloss](https://github.com/charmbracelet/lipgloss)** — TUI styling
- **YAML** — configuration format (LLM-friendly for AI-assisted editing)

## Requirements

- Go 1.21+

## Installation

### From source

```bash
git clone https://github.com/dan-westall/cli-chameleon.git
cd chameleon
make install
```

This installs the `chameleon` binary to your `$GOPATH/bin`.

### Manual build

```bash
make build
cp bin/chameleon /usr/local/bin/
```

### Go install

```bash
go install github.com/dan-westall/cli-chameleon@latest
```

## Usage

Navigate to any project directory and run:

```bash
chameleon
```

### First run (no config)

If no `chameleon.yaml` exists in the current directory, Chameleon creates a commented template:

```bash
$ cd my-new-project
$ chameleon
Created chameleon.yaml — edit it to configure your commands.
```

### With config

Once configured, running `chameleon` launches the interactive TUI menu showing all available commands for that project.

### Version

```bash
chameleon --version
```

## Configuration

Create a `chameleon.yaml` in your project root:

```yaml
name: "My Project"
commands:
  - name: ingest
    description: "Run the ingestion pipeline"
    run: "python scripts/ingest.py"
    stream: true

  - name: build
    description: "Build the project"
    run:
      - "npm install"
      - "npm run build"

  - name: deploy
    description: "Deploy services in parallel"
    run:
      - "docker push app"
      - "docker push worker"
    parallel: true

  - name: test
    description: "Run test suite"
    run: "go test ./..."
```

### Fields

| Field | Type | Default | Description |
|-------|------|---------|-------------|
| `name` | string | required | Command name displayed in the menu |
| `description` | string | `""` | Short description shown alongside the name |
| `run` | string or string[] | required | Shell command(s) to execute |
| `stream` | bool | `false` | Show live output in a split-panel view |
| `parallel` | bool | `false` | Run array commands concurrently instead of sequentially |

### Execution Modes

**Single command** — `run` is a string:

```yaml
run: "npm run build"
```

**Sequential** — `run` is an array (default behaviour):

```yaml
run:
  - "npm install"
  - "npm run build"
```

Commands run in order. Execution stops on the first failure.

**Parallel** — `run` is an array with `parallel: true`:

```yaml
run:
  - "docker push app"
  - "docker push worker"
parallel: true
```

All commands run concurrently. Fails if any command fails.

### Stream vs Non-Stream

**Non-stream commands** (default) execute in the background. On completion, a success or failure toast flashes briefly and the menu returns.

**Stream commands** (`stream: true`) open a 70/30 split-panel view:

```
┌──────────────────────────────────────────────┬──────────────────┐
│ Processing batch 1/10...                     │ ingest           │
│ Processing batch 2/10...                     │                  │
│ Processing batch 3/10...                     │ ● Running        │
│ Processing batch 4/10...                     │                  │
│                                              │                  │
├──────────────────────────────────────────────┴──────────────────┤
│ esc back • q quit                                               │
└─────────────────────────────────────────────────────────────────┘
```

- **70% left panel** — live-streamed stdout/stderr from the running process
- **30% right panel** — command name and status (reserved for future use)

## Controls

| Key | Action |
|-----|--------|
| `↑` / `k` | Move cursor up |
| `↓` / `j` | Move cursor down |
| `Enter` | Execute selected command |
| `Esc` | Return to menu (from stream view) |
| `q` | Quit |

## Project Structure

```
chameleon/
├── main.go                        # Entry point
├── Makefile                       # Build targets
├── go.mod
└── internal/
    ├── config/
    │   ├── config.go              # YAML schema and parser
    │   └── template.go            # Template generation
    ├── executor/
    │   └── executor.go            # Command execution (sequential/parallel/stream)
    └── tui/
        ├── model.go               # Root Bubble Tea model
        ├── menu.go                # Main menu view
        ├── toast.go               # Success/fail toast
        └── stream.go              # 70/30 split stream view
```

## Build

```bash
make build      # Build binary to bin/chameleon
make install    # Install to $GOPATH/bin
make test       # Run all tests
make clean      # Remove build artefacts
```

## Example Configurations

### Node.js project

```yaml
name: "frontend-app"
commands:
  - name: dev
    description: "Start dev server"
    run: "npm run dev"
    stream: true
  - name: build
    description: "Production build"
    run:
      - "npm ci"
      - "npm run build"
  - name: test
    description: "Run tests"
    run: "npm test"
  - name: lint
    description: "Lint and fix"
    run: "npm run lint:fix"
```

### Python data pipeline

```yaml
name: "data-pipeline"
commands:
  - name: ingest
    description: "Run full ingestion cycle"
    run: "python -m pipeline.ingest"
    stream: true
  - name: validate
    description: "Validate data quality"
    run: "python -m pipeline.validate"
    stream: true
  - name: setup
    description: "Install dependencies"
    run:
      - "python -m venv .venv"
      - ".venv/bin/pip install -r requirements.txt"
```

### Infrastructure project

```yaml
name: "platform-infra"
commands:
  - name: plan
    description: "Terraform plan"
    run: "terraform plan"
    stream: true
  - name: apply
    description: "Terraform apply"
    run: "terraform apply -auto-approve"
    stream: true
  - name: build
    description: "Build all containers"
    run:
      - "docker build -t app ./app"
      - "docker build -t worker ./worker"
    parallel: true
```

## Contributing

Contributions are welcome. Please open an issue to discuss changes before submitting a pull request.

## Licence

MIT
