package manager

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"github.com/google/uuid"
	"github.com/kavos113/seseragi/model"
	"github.com/kavos113/seseragi/model/repository/json"
)

var (
	jsonRepo     *json.JsonRepository
	taskRepo     model.TaskRepository
	workflowRepo model.WorkflowRepository
	repoInitOnce sync.Once

	dockerClient *DockerClient
)

func InitRepository() {
	repoInitOnce.Do(func() {
		appDataDir, err := os.UserConfigDir()
		if err != nil {
			panic(fmt.Sprintf("Failed to get user config directory: %v", err))
		}

		jsonRepo = json.NewJsonRepository(filepath.Join(appDataDir, "seseragi"))
		taskRepo = json.NewJSONTaskRepository(jsonRepo)
		workflowRepo = json.NewJSONWorkflowRepository(jsonRepo)

		dockerClient = NewDockerClient()
	})
}

// yamlPath is expected to be a ABSOLUTE PATH
func BuildTask(yamlPath string) error {
	yamlData, err := os.ReadFile(yamlPath)
	if err != nil {
		return err
	}

	taskInfo, err := LoadTaskInfoFromYAML(yamlData, yamlPath)
	if err != nil {
		return err
	}

	id := uuid.New().String()
	task := model.Task{
		ID:        id,
		Name:      taskInfo.Name,
		ImageName: fmt.Sprintf("%s-%s", taskInfo.Name, id),
		YamlPath:  yamlPath,
	}

	if err := dockerClient.BuildImage(filepath.Dir(yamlPath), task.ImageName); err != nil {
		return err
	}

	if _, err := taskRepo.CreateTask(task); err != nil {
		return err
	}

	return nil
}

func AddWorkflow(yamlPath string) error {
	yamlData, err := os.ReadFile(yamlPath)
	if err != nil {
		return err
	}

	workflowInfo, err := LoadWorkflowInfoFromYAML(yamlData, yamlPath)
	if err != nil {
		return err
	}

	workflow, err := ParseWorkflow(workflowInfo, yamlPath)
	if err != nil {
		return err
	}

	if _, err := workflowRepo.CreateWorkflow(workflow); err != nil {
		return err
	}

	return nil
}

func ListWorkflows() ([]model.Workflow, error) {
	return workflowRepo.GetAllWorkflows()
}

func ListTasks() ([]model.Task, error) {
	return taskRepo.GetAllTasks()
}

func main() {}
