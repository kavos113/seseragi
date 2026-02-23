package main

import (
	"os"

	"github.com/kavos113/seseragi/manager-client/cli/commands"
)

func main() {
	if len(os.Args) < 2 {
		println("Usage: seseragi-cli <command> [args]")
		os.Exit(1)
	}

	if os.Args[1] == "build" {
		if len(os.Args) < 3 {
			println("Usage: seseragi-cli build <yaml_path>")
			os.Exit(1)
		}
		if err := commands.BuildCommand(os.Args[2]); err != nil {
			println("Error:", err.Error())
			os.Exit(1)
		}
		println("Task built successfully")
	}
}
