package manager

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/google/uuid"
	"github.com/kavos113/seseragi/model"
	"github.com/kavos113/seseragi/model/repository/json"
)

var (
	jsonRepo     = json.NewJsonRepository("data")
	taskRepo     = json.NewJSONTaskRepository(jsonRepo)
	workflowRepo = json.NewJSONWorkflowRepository(jsonRepo)
	dockerClient = NewDockerClient()
)

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

	nodes := make([]model.Node, 0, len(workflowInfo.Nodes))
	for taskName, nodeInfo := range workflowInfo.Nodes {
		task, err := taskRepo.GetTaskByID(nodeInfo.ID)
		if err != nil {
			return fmt.Errorf("failed to get task by ID %s for node %s: %w", nodeInfo.ID, taskName, err)
		}
		nodes = append(nodes, model.Node{
			TaskID:       task.ID,
			Dependencies: []model.Task{},
		})
	}

	id := uuid.New().String()
	workflow := model.Workflow{
		ID:          id,
		Name:        workflowInfo.Name,
		Nodes:       nodes,
		YamlPath:    yamlPath,
	}	

	if _, err := workflowRepo.CreateWorkflow(workflow); err != nil {
		return err
	}

	return nil
}

func main() {}
