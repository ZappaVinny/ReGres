package handlers

import (
	h "regres/cli/internal/helpers"
)

func SrvShell() {
	SrvCommand("sh")
}

func SrvCommand(args ...string) {
	cmd := append([]string{"exec", "srv"}, args...)
	SrvUp(true)
	h.RunCompose(cmd...)
	SrvDown(true)
}

func SrvUp(quiet bool) {
	h.RunComposeQuiet(quiet, "up", "-d", "srv")
}

func SrvDown(quiet bool) {
	h.RunComposeQuiet(quiet, "stop", "srv")
}

const SrvHelpText = `
rgs srv runs a command in the srv container.

Usage:
  rgs srv <command> [args]

Available Commands:
  help        Show help for srv command
  sh          Start a shell in the srv container
  up		  Start the srv container
  down		  Stop the srv container

Any other command is passed through and run directly inside the srv
container, e.g.:

  rgs srv go test ./...
  rgs srv go mod tidy
  rgs srv sh`
