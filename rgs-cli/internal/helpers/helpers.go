package helpers

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func RunCompose(args ...string) {
	RunComposeQuiet(false, args...)
}

func RunComposeQuiet(quiet bool, args ...string) {
	cmd := exec.Command("docker", append([]string{"compose"}, args...)...)
	cmd.Env = os.Environ()
	cmd.Stdin = os.Stdin
	if !quiet {
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
	}

	if err := cmd.Run(); err != nil {
		fmt.Printf("Error running command: %v\n", err)
	}
}

func RunMigrationsUp() {
	RunCompose("run", "--rm", "srv", "sh", "-c", `migrate -path db/migrations -database "$DATABASE_URL" up`)
}

func RunMigrationsDown() {
	RunCompose("run", "--rm", "srv", "sh", "-c", `migrate -path db/migrations -database "$DATABASE_URL" down -all`)
}

func RunMigrationsReset() {
	script := `
		migrate -path db/migrations -database "$DATABASE_URL" down -all || true
		migrate -path db/migrations -database "$DATABASE_URL" up
	`
	RunCompose("run", "--rm", "srv", "sh", "-c", script)
}

// Not actually sure I need this, I shifted from what I was doing, but ill keep it here
func GetActiveContainers() ([]string, error) {
	cmd := exec.Command("docker", "compose", "ps", "--format", "{{.Name}}")
	cmd.Env = os.Environ()

	out, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("error running command: %w", err)
	}

	var names []string
	for _, line := range strings.Split(strings.TrimSpace(string(out)), "\n") {
		line = strings.TrimSpace(line)
		if line != "" {
			names = append(names, line)
		}
	}

	return names, nil
}
