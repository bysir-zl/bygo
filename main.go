package main

import (
    "os"
    "github.com/bysir-zl/bygo/artisan"
)

func main() {
    args := os.Args
    args = []string{"create", "helloworld"}

    command := args[0]

    switch command {
    case "create":
        projectName := args[1]
        artisan.CreateProject(projectName)

    }
}
