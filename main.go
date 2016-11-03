package main

import (
	"github.com/bysir-zl/bygo/artisan"
	"os"
)

func main() {
	args := os.Args
	if len(args) == 1 {
		args = []string{"create", "helloworld"}
	}

	command := args[1]

	switch command {
	case "create":
		projectName := args[2]
		artisan.CreateProject(projectName)

	}
}
