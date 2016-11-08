package main

import (
	"github.com/bysir-zl/bygo/artisan"
	"github.com/deepzz0/go-com/log"
	"os"
)

func main() {
	args := os.Args
	//if len(args) == 1 {
	//	args = []string{"create", "helloworld"}
	//}

	if len(args) == 1 {
		log.Print("error args")
		return
	}

	command := args[1]

	switch command {
	case "create":
		projectName := args[2]
		artisan.CreateProject(projectName)
	case "model":
		table := args[2]
		artisan.CreateModelFile(table)
	case "swagger":
		if len(args) != 4 {
			args = []string{"", "", "./", "./swagger.json"}
		}
		path := args[2]
		out := args[3]
		artisan.Swagger(path, out)
	}
}
