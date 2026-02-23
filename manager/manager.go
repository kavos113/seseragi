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

func main() {}
