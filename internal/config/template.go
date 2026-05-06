package config

import "os"

const template = `# Chameleon CLI Configuration
# This file defines the commands available in your project's TUI menu.
#
# Each command has:
#   name:        Short name displayed in the menu
#   description: What the command does
#   run:         Command to execute (string or array of strings)
#   stream:      Show live output in split view (default: false)
#   parallel:    Run array commands concurrently (default: false)
#
# Examples:
#
# commands:
#   - name: build
#     description: "Build the project"
#     run: "npm run build"
#
#   - name: deploy
#     description: "Build and deploy"
#     run:
#       - "npm run build"
#       - "npm run deploy"
#
#   - name: ingest
#     description: "Run ingestion pipeline"
#     run: "python scripts/ingest.py"
#     stream: true

name: ""
commands: []
`

func CreateTemplate() error {
	return os.WriteFile(FileName, []byte(template), 0644)
}
