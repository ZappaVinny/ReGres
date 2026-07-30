package main

import (
	"fmt"
	"os"

	c "regres/cli/internal/handlers"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println(helpText)
		os.Exit(1)
	}

	switch os.Args[1] {
	case "help":
		fmt.Println(helpText)
	case "init":
		c.EnvInit()
	case "run":
		c.EnvRun()
	case "down":
		c.EnvDown()

	case "srv":
		if len(os.Args) < 3 {
			fmt.Println(c.SrvHelpText)
			os.Exit(1)
		}
		switch os.Args[2] {
		case "help":
			fmt.Println(c.SrvHelpText)
		case "sh":
			c.SrvShell()
		case "up":
			c.SrvUp(false)
		case "down":
			c.SrvDown(false)
		default:
			c.SrvCommand(os.Args[2:]...)
		}

	case "web":
		if len(os.Args) < 3 {
			fmt.Println(c.WebHelpText)
			os.Exit(1)
		}
		switch os.Args[2] {
		case "help":
			fmt.Println(c.WebHelpText)
		case "sh":
			c.WebShell()
		case "up":
			c.WebUp(false)
		case "down":
			c.WebDown(false)
		default:
			c.WebCommand(os.Args[2:]...)
		}

	case "db":
		if len(os.Args) < 3 {
			fmt.Println(c.DBHelpText)
			os.Exit(1)
		}
		switch os.Args[2] {
		case "help":
			fmt.Println(c.DBHelpText)
		case "sh":
			c.DBShell()
		case "up":
			c.DBUp(false)
		case "down":
			c.DBDown(false)
		case "migrate":
			if len(os.Args) < 4 {
				fmt.Println(c.DBMigrateHelpText)
				os.Exit(1)
			}
			switch os.Args[3] {
			case "help":
				fmt.Println(c.DBMigrateHelpText)
			case "up":
				c.DBMigrateUp()
			case "down":
				c.DBMigrateDown()
			case "reset":
				c.DBMigrateReset()
			default:
				fmt.Println(c.DBMigrateHelpText)
				os.Exit(1)
			}
		default:
			c.DBCommand(os.Args[2:]...)
		}

	default:
		fmt.Printf("unknown command: %s\n", os.Args[1])
		os.Exit(1)
	}
}

const helpText = `
ReGres CLI is a command-line tool for managing the ReGres development environment.

Usage:
  rgs <command> [options]

Available Commands:
  init        Initialize a ReGres development environment
  run         Run a ReGres development environment
  down        Stop and remove the ReGres environment
  srv         Run a command in the srv container
  web         Run a command in the web container
  db          Run a command in the db container, incl. migrations
  help        Show help for a command

Use "rgs help <command>" for more information about a command.`
