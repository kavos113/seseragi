package main

import (
	"fmt"
	"os"
	"time"

	"github.com/google/uuid"
	"github.com/kavos113/seseragi/model"
	"github.com/kavos113/seseragi/runner/service"
)

func main() {
	if len(os.Args) < 2 {
		println("Usage: seseragi [workflow-id]")
	}

	dc, err := service.NewDockerClient()
	if err != nil {
		panic(err)
	}

	wr := service.NewWorkflowRunner()

	workflow, err := wr.GetWorkflowByID(os.Args[1])
	if err != nil {
		panic(err)
	}

	fmt.Printf("Running workflow: %s\n", workflow.Name)
	id := uuid.New().String()
	run := model.WorkflowRun{
		ID:         id,
		WorkflowID: workflow.ID,
		StartTime:  time.Now(),
	}

	err = wr.RunWorkflow(workflow, func(n model.Node) error {
		imageName, err := wr.GetImageNameByTaskID(n.TaskID)
		if err != nil {
			return err
		}

		return dc.RunContainer(imageName)
	})
	run.EndTime = time.Now()
	if err != nil {
		run.Status = model.WorkflowStatusFailed
		fmt.Printf("Workflow failed: %s\n", err.Error())
	} else {
		run.Status = model.WorkflowStatusCompleted
	}

	if err := wr.SaveWorkflowRun(run); err != nil {
		panic(err)
	}

	dc.Cleanup()
}
