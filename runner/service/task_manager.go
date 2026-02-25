package service

import (
	"os"
	"path/filepath"

	"github.com/kavos113/seseragi/model"
	"github.com/kavos113/seseragi/model/repository/json"
)

type TaskManager struct {
	repo model.TaskRepository
}

func NewTaskManager(repo model.TaskRepository) *TaskManager {
	userConfigDir, err := os.UserConfigDir()
	if err != nil {
		panic(err)
	}

	jsonRepo := json.NewJsonRepository(filepath.Join(userConfigDir, "seseragi"))
	return &TaskManager{
		repo: json.NewJSONTaskRepository(jsonRepo),
	}
}

func (tm *TaskManager) GetImageNameByTaskID(taskID string) (string, error) {
	task, err := tm.repo.GetTaskByID(taskID)
	if err != nil {
		return "", err
	}
	return task.ImageName, nil
}
