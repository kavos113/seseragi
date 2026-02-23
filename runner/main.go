package main

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/kavos113/seseragi/model"
	"github.com/kavos113/seseragi/runner/service"
)

func main() {
	dc, err := service.NewDockerClient()
	if err != nil {
		panic(err)
	}

	wr := service.NewWorkflowRunner()

	workflows, err := wr.GetWorkflowToRun()
	if err != nil {
		panic(err)
	}

	for _, workflow := range workflows {
		start := time.Now()
		fmt.Printf("Running workflow: %s\n", workflow.Name)

		for _, node := range workflow.Nodes {
			imageName, err := wr.GetImageNameByTaskID(node.TaskID)
			if err != nil {
				panic(err)
			}

			if err := dc.RunContainer(imageName); err != nil {
				panic(err)
			}
		}

		id := uuid.New().String()
		run := model.WorkflowRun{
			ID:         id,
			WorkflowID: workflow.ID,
			StartTime:  start,
			EndTime:    time.Now(),
			Status:     model.WorkflowStatusCompleted,
		}
		if err := wr.SaveWorkflowRun(run); err != nil {
			panic(err)
		}
	}

	defer dc.Cleanup()
}
