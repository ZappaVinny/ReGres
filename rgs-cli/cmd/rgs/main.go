package main

import (
	"fmt"
	"os"

	cmd "regres/cli/internal/handlers"
	h "regres/cli/internal/helpers"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println(h.HelpText)
		os.Exit(1)
	}

	switch os.Args[1] {
	case "help":
		fmt.Println(h.HelpText)
	case "init":
		cmd.RunInit()
	case "run":
		cmd.RunDev()
	case "down":
		cmd.RunDown()
	default:
		fmt.Printf("unknown command: %s\n", os.Args[1])
		os.Exit(1)
	}
}
