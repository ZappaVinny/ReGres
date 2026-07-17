package handlers

import (
	"fmt"
	h "regres/cli/internal/helpers"
)

func RunInit() {
	h.RunCompose("build")
	fmt.Println("Initializing database+++")
	h.RunCompose("up", "-d", "db")

	fmt.Println("Running migrations+++")
	RunMigrationsReset()

	fmt.Println("Constructing Models+++")
	h.RunCompose("run", "--rm", "srv", "sqlc", "generate")

	fmt.Println("Constructing Models+++")
	h.RunCompose("create", "web", "srv")

	fmt.Println(h.InitText)

}

func RunDown() {
	h.RunCompose("down")
}

func RunDev() {
	h.RunCompose("up", "web", "srv")
}

func RunMigrationsReset() {
	script := `
		migrate -path db/migrations -database "$DATABASE_URL" down -all || true
		migrate -path db/migrations -database "$DATABASE_URL" up
	`
	h.RunCompose("run", "--rm", "srv", "sh", "-c", script)
}
