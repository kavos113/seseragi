package domain

type Task struct {
	ID       string
	Name     string // unique
	YamlPath string
	TaskDef  TaskDefinition
}

type TaskType string

const (
	TaskTypeContainer TaskType = "container"
	TaskTypeCommand   TaskType = "command"
)

type TaskDefinition interface {
	Type() TaskType
}

type ContainerTaskDefinition struct {
	ImageName string
}

func (c ContainerTaskDefinition) Type() TaskType {
	return TaskTypeContainer
}

type CommandTaskDefinition struct {
	Command    string
	WorkingDir string
}

func (c CommandTaskDefinition) Type() TaskType {
	return TaskTypeCommand
}

type TaskRepository interface {
	CreateTask(task Task) (Task, error)
	GetTaskByID(id string) (Task, error)
	GetTaskByName(name string) (Task, error)
	GetAllTasks() ([]Task, error)

	UpdateTask(task Task) (Task, error)
	DeleteTask(id string) error
}
