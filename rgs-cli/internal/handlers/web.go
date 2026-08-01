package handlers

import (
	h "regres/cli/internal/helpers"
)

func WebShell() {
	WebCommand("sh")
}

func WebCommand(args ...string) {
	cmd := append([]string{"exec", "web"}, args...)
	WebUp(true)
	h.RunCompose(cmd...)
	WebDown(true)
}

func WebUp(quiet bool) {
	h.RunComposeQuiet(quiet, "up", "-d", "web")
}

func WebDown(quiet bool) {
	h.RunComposeQuiet(quiet, "stop", "web")
}

const WebHelpText = `
rgs web runs a command in the web container.

Usage:
  rgs web <command> [args]

Available Commands:
  help        Show help for web command
  sh          Start a shell in the web container
  up		  Start the web container
  down		  Stop the web container

Any other command is passed through and run directly inside the web
container, e.g.:

  rgs web npm install
  rgs web npm run build
  rgs web sh`
