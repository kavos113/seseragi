package commands

import (
	"fmt"

	"github.com/kavos113/seseragi/manager"
)

func ListTasks() error {
	tasks, err := manager.ListTasks()
	if err != nil {
		return err
	}

	for _, task := range tasks {
		fmt.Printf("%s - %s - %s\n", task.ID, task.Name, task.YamlPath)
	}

	return nil
}
