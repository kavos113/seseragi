package main

import (
	"os"

	"github.com/kavos113/seseragi/manager"
	"github.com/kavos113/seseragi/manager-client/cli/commands"
)

func main() {
	if len(os.Args) < 2 {
		println("Usage: seseragi-cli <command> [args]")
		os.Exit(1)
	}

	manager.InitRepository()

	switch os.Args[1] {
	case "build":
		if len(os.Args) < 3 {
			println("Usage: seseragi-cli build <yaml_path>")
			os.Exit(1)
		}
		if err := commands.BuildCommand(os.Args[2]); err != nil {
			println("Error:", err.Error())
			os.Exit(1)
		}
		println("Task built successfully")
	case "add-workflow":
		if len(os.Args) < 3 {
			println("Usage: seseragi-cli add-workflow <yaml_path>")
			os.Exit(1)
		}
		if err := commands.AddWorkflow(os.Args[2]); err != nil {
			println("Error:", err.Error())
			os.Exit(1)
		}
		println("Workflow added successfully")
	default:
		println("Unknown command:", os.Args[1])
		os.Exit(1)
	}
}
