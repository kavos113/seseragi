package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
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

	check := func() {
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
	}

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	ticker := time.NewTicker(5 * time.Minute)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			check()

		case <-ctx.Done():
			fmt.Println("Received shutdown signal")
			dc.Cleanup()
			return
		}
	}
}
