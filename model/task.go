package model

type Task struct {
	ID        string
	Name      string
	ImageName string
}

type TaskRepository interface {
	CreateTask(task Task) (Task, error)
	GetTaskByID(id string) (Task, error)

	UpdateTask(task Task) (Task, error)
	DeleteTask(id string) error
}
