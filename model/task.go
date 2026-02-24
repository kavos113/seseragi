package model

type Task struct {
	ID        string
	Name      string // unique
	ImageName string
	YamlPath  string
}

type TaskRepository interface {
	CreateTask(task Task) (Task, error)
	GetTaskByID(id string) (Task, error)
	GetTaskByName(name string) (Task, error)
	GetAllTasks() ([]Task, error)

	UpdateTask(task Task) (Task, error)
	DeleteTask(id string) error
}
