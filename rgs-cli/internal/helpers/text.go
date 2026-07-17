package helpers

const HelpText = `
ReGres CLI is a command-line tool for managing the ReGres development environment.

Usage:
  rgs <command> [options]

Available Commands:
  init        Initialize a ReGres development environment
  run         Run a ReGres development environment
  help        Show help for a command

Use "rgs help <command>" for more information about a command.`

const SrvHelpText = `
rgs srv runs a command in the srv container.

Usage:
  rgs srv <command> [args]

Available Commands:
  help        Show help for srv command

Any other command is passed through and run directly inside the srv
container, e.g.:

  rgs srv go test ./...
  rgs srv go mod tidy
  rgs srv sh`

const WebHelpText = `
rgs web runs a command in the web container.

Usage:
  rgs web <command> [args]

Available Commands:
  help        Show help for web command

Any other command is passed through and run directly inside the web
container, e.g.:

  rgs web npm install
  rgs web npm run build
  rgs web sh`

const InitText = `
ReGres initialized.
Run the stack with:

rgs run`
