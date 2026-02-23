package main

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/kavos113/seseragi/model"
)

func main() {
	dc, err := NewDockerClient()
	if err != nil {
		panic(err)
	}

	wr := NewWorkflowRunner()

	workflows, err := wr.GetWorkflowToRun()
	if err != nil {
		panic(err)
	}

	for _, workflow := range workflows {
		start := time.Now()
		fmt.Printf("Running workflow: %s\n", workflow.Name)
		
		for _, node := range workflow.Nodes {
			if err := dc.RunContainer(node.TaskID); err != nil {
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
