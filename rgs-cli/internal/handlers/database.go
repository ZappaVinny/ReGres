package handlers

import (
	"fmt"

	h "regres/cli/internal/helpers"
)

func DBShell() {
	DBCommand("sh")
}

func DBCommand(args ...string) {
	cmd := append([]string{"exec", "db"}, args...)
	DBUp(true)
	h.RunCompose(cmd...)
}

func DBUp(quiet bool) {
	h.RunComposeQuiet(quiet, "up", "-d", "db")
}

func DBDown(quiet bool) {
	h.RunComposeQuiet(quiet, "stop", "db")
}

func DBMigrateUp() {
	fmt.Println("Running migrations up+++")
	h.RunCompose("up", "-d", "db")
	h.RunMigrationsUp()
}

func DBMigrateDown() {
	fmt.Println("Running migrations down+++")
	h.RunCompose("up", "-d", "db")
	h.RunMigrationsDown()
}

func DBMigrateReset() {
	fmt.Println("Resetting database+++")
	h.RunCompose("up", "-d", "db")
	h.RunMigrationsReset()

	fmt.Println("Constructing Models+++")
	h.RunCompose("run", "--rm", "srv", "sqlc", "generate")
}

const DBHelpText = `
rgs db runs a command in the db container.

Usage:
  rgs db <command> [args]

Available Commands:
  help        Show help for db command
  sh          Start a shell in the db container
  up          Start the db container
  down        Stop the db container
  migrate     Run database migrations (see "rgs db migrate help")

Any other command is passed through and run directly inside the db
container, e.g.:

  rgs db psql -U postgres -d regres`

const DBMigrateHelpText = `
rgs db migrate runs go-migrate and sqlc generate against the db container.

Usage:
  rgs db migrate <command>

Available Commands:
  help        Show help for db migrate command
  up          Run all pending migrations
  down        Roll back all migrations
  reset       Roll back and re-run all migrations, then regenerate sqlc models`
