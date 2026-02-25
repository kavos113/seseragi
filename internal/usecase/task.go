package usecase

import (
	"github.com/kavos113/seseragi/internal/domain"
)

type TaskUseCase interface {
	AddTask(task domain.Task) error
	ListTasks() ([]domain.Task, error)
	DeleteTask(taskID string) error
}

type taskUseCase struct {
	taskRepo     domain.TaskRepository
	taskProvider domain.TaskProvider
	idGenerator  IDGenerator
}

func NewTaskUseCase(taskRepo domain.TaskRepository, taskProvider domain.TaskProvider, idGenerator IDGenerator) TaskUseCase {
	return &taskUseCase{
		taskRepo:     taskRepo,
		taskProvider: taskProvider,
		idGenerator:  idGenerator,
	}
}

func (uc *taskUseCase) AddTask(task domain.Task) error {
	id := uc.idGenerator.GenerateID()
	task.ID = id
	if err := uc.taskProvider.BuildTask(task); err != nil {
		return err
	}

	_, err := uc.taskRepo.CreateTask(task)
	return err
}

func (uc *taskUseCase) ListTasks() ([]domain.Task, error) {
	return uc.taskRepo.GetAllTasks()
}

func (uc *taskUseCase) DeleteTask(taskID string) error {
	return uc.taskRepo.DeleteTask(taskID)
}
