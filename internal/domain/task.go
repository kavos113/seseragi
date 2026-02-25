package domain

type Task struct {
	ID       string
	Name     string // unique
	YamlPath string
	TaskDef  TaskDefinition
}

type TaskType string

const (
	TaskTypeDocker  TaskType = "docker"
	TaskTypeCommand TaskType = "command"
)

type TaskDefinition interface {
	Type() TaskType
}

type DockerTaskDefinition struct {
	ImageName  string
	ContextDir string
}

func (c DockerTaskDefinition) Type() TaskType {
	return TaskTypeDocker
}

type CommandTaskDefinition struct {
	Command    string
	WorkingDir string
}

func (c CommandTaskDefinition) Type() TaskType {
	return TaskTypeCommand
}

// taskの永続化
type TaskRepository interface {
	CreateTask(task Task) (Task, error)
	GetTaskByID(id string) (Task, error)
	GetTaskByName(name string) (Task, error)
	GetAllTasks() ([]Task, error)

	UpdateTask(task Task) (Task, error)
	DeleteTask(id string) error
}

// コンテナのビルドやデータの準備など
type TaskProvider interface {
	BuildTask(task Task) error
}
