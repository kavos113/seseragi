package usecase

import (
	"github.com/kavos113/seseragi/internal/domain"
)

type TaskUseCase interface {
	AddTask(task domain.Task, provider domain.TaskProvider) error
	UpdateTask(task domain.Task, provider domain.TaskProvider) error
	ListTasks() ([]domain.Task, error)
	DeleteTask(taskID string) error
}

type taskUseCase struct {
	taskRepo    domain.TaskRepository
	idGenerator IDGenerator
}

func NewTaskUseCase(taskRepo domain.TaskRepository, idGenerator IDGenerator) TaskUseCase {
	return &taskUseCase{
		taskRepo:    taskRepo,
		idGenerator: idGenerator,
	}
}

func (uc *taskUseCase) AddTask(task domain.Task, provider domain.TaskProvider) error {
	id := uc.idGenerator.GenerateID()
	task.ID = id

	if dockerDef, ok := task.TaskDef.(domain.DockerTaskDefinition); ok {
		task.TaskDef = domain.DockerTaskDefinition{
			ImageName:  id,
			ContextDir: dockerDef.ContextDir,
		}
	}

	if err := provider.BuildTask(task); err != nil {
		return err
	}

	_, err := uc.taskRepo.CreateTask(task)
	return err
}

func (uc *taskUseCase) UpdateTask(task domain.Task, provider domain.TaskProvider) error {
	existingTask, err := uc.taskRepo.GetTaskByID(task.ID)
	if err != nil {
		return err
	}

	if dockerDef, ok := task.TaskDef.(domain.DockerTaskDefinition); ok {
		task.TaskDef = domain.DockerTaskDefinition{
			ImageName:  existingTask.ID,
			ContextDir: dockerDef.ContextDir,
		}
	}

	if err := provider.BuildTask(task); err != nil {
		return err
	}
	
	_, err = uc.taskRepo.UpdateTask(task)
	return err
}

func (uc *taskUseCase) ListTasks() ([]domain.Task, error) {
	return uc.taskRepo.GetAllTasks()
}

func (uc *taskUseCase) DeleteTask(taskID string) error {
	return uc.taskRepo.DeleteTask(taskID)
}
