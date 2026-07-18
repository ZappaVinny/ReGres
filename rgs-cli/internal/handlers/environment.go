package handlers

import (
	"fmt"
	h "regres/cli/internal/helpers"
)

func EnvInit() {
	h.RunCompose("build")
	fmt.Println("Initializing database+++")
	h.RunCompose("up", "-d", "db")

	fmt.Println("Running migrations+++")
	h.RunMigrationsReset()

	fmt.Println("Constructing Models+++")
	h.RunCompose("run", "--rm", "srv", "sqlc", "generate")

	fmt.Println("Constructing Models+++")
	h.RunCompose("create", "web", "srv")

	fmt.Println(initText)

}

func EnvDown() {
	h.RunCompose("down")
	fmt.Println("Environment down+++")
}

func EnvRun() {
	h.RunCompose("up", "web", "srv")
}

const initText = `
ReGres initialized.
Run the stack with:

rgs run`
