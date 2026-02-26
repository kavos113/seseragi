package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"

	"github.com/kavos113/seseragi/internal/adapter/docker"
	"github.com/kavos113/seseragi/internal/domain"
	"github.com/kavos113/seseragi/internal/repository/json"
	"github.com/kavos113/seseragi/internal/runner/command"
	dockerr "github.com/kavos113/seseragi/internal/runner/docker"
	"github.com/kavos113/seseragi/internal/usecase"
)

func main() {
	dc, err := docker.NewClient()
	if err != nil {
		panic(err)
	}

	appDataDir, err := os.UserConfigDir()
	if err != nil {
		panic(err)
	}
	jsonRepo := json.NewJsonRepository(filepath.Join(appDataDir, "seseragi"))
	taskRepo := json.NewJSONTaskRepository(jsonRepo)
	workflowRepo := json.NewJSONWorkflowRepository(jsonRepo)
	workflowRunRepo := json.NewJSONWorkflowRunRepository(jsonRepo)

	idGenerator := usecase.NewUUIDGenerator()

	dr := dockerr.NewDockerNodeRunner(dc)

	wu := usecase.NewWorkflowUseCase(workflowRepo, taskRepo, idGenerator)
	wru := usecase.NewWorkflowRunUseCase(workflowRepo, workflowRunRepo, taskRepo, idGenerator)

	runnerSelector := func(node domain.Node) domain.NodeRunner {
		t, err := wu.GetTaskTypeFromNode(node)
		if err != nil {
			fmt.Printf("Error determining task type for node %s: %v\n", node.Name, err)
			return nil
		}

		switch t {
		case domain.TaskTypeDocker:
			return dr

		case domain.TaskTypeCommand:
			return command.NewCommandTaskRunner(5 * time.Minute)

		default:
			fmt.Printf("No runner available for task type %s in node %s\n", t, node.Name)
			return nil
		}
	}

	check := func() {
		toRun, err := wru.GetWorkflowsToRun()
		if err != nil {
			fmt.Printf("Error getting workflows to run: %v\n", err)
			return
		}

		for _, wf := range toRun {
			go func(wf domain.Workflow) {
				fmt.Printf("Starting workflow run for workflow %s\n", wf.ID)

				err := wru.RunWorkflow(wf.ID, runnerSelector)
				if err != nil {
					fmt.Printf("Error running workflow %s: %v\n", wf.ID, err)
				}
			}(wf)
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
			fmt.Println("Shutting down seseragi daemon...")
			return
		}
	}
}
