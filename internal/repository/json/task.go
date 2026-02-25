package json

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"slices"

	"github.com/kavos113/seseragi/internal/domain"
)

type taskDTO struct {
	ID       string          `json:"id"`
	Name     string          `json:"name"`
	YamlPath string          `json:"yaml_path"`
	Type     domain.TaskType `json:"type"`
	TaskDef  json.RawMessage `json:"task_def"`
}

func (t taskDTO) toDomain() (domain.Task, error) {
	if t.Type == "" {
		return domain.Task{
			ID:       t.ID,
			Name:     t.Name,
			YamlPath: t.YamlPath,
			TaskDef:  nil,
		}, nil
	}

	var taskDef domain.TaskDefinition
	switch t.Type {
	case domain.TaskTypeDocker:
		var dockerDef dockerTaskDefDTO
		err := json.Unmarshal(t.TaskDef, &dockerDef)
		if err != nil {
			return domain.Task{}, err
		}
		taskDef = dockerDef.toDomain()

	case domain.TaskTypeCommand:
		var commandDef commandTaskDefDTO
		err := json.Unmarshal(t.TaskDef, &commandDef)
		if err != nil {
			return domain.Task{}, err
		}
		taskDef = commandDef.toDomain()

	default:
		return domain.Task{}, errors.New("invalid task type")
	}

	return domain.Task{
		ID:       t.ID,
		Name:     t.Name,
		YamlPath: t.YamlPath,
		TaskDef:  taskDef,
	}, nil
}

func fromDomainTask(task domain.Task) (taskDTO, error) {
	if task.TaskDef == nil {
		return taskDTO{
			ID:       task.ID,
			Name:     task.Name,
			YamlPath: task.YamlPath,
			Type:     "",
			TaskDef:  nil,
		}, nil
	}

	var taskDefDTO any
	switch task.TaskDef.Type() {
	case domain.TaskTypeDocker:
		dockerDef, ok := task.TaskDef.(domain.DockerTaskDefinition)
		if !ok {
			return taskDTO{}, errors.New("invalid docker task definition")
		}
		taskDefDTO, _ = fromDomainDockerTaskDef(dockerDef)

	case domain.TaskTypeCommand:
		commandDef, ok := task.TaskDef.(domain.CommandTaskDefinition)
		if !ok {
			return taskDTO{}, errors.New("invalid command task definition")
		}
		taskDefDTO, _ = fromDomainCommandTaskDef(commandDef)

	default:
		return taskDTO{}, errors.New("invalid task type")
	}

	taskDefData, err := json.Marshal(taskDefDTO)
	if err != nil {
		return taskDTO{}, err
	}

	return taskDTO{
		ID:       task.ID,
		Name:     task.Name,
		YamlPath: task.YamlPath,
		Type:     task.TaskDef.Type(),
		TaskDef:  taskDefData,
	}, nil
}

type dockerTaskDefDTO struct {
	ImageName  string `json:"image_name"`
	ContextDir string `json:"context_dir"`
}

func (d dockerTaskDefDTO) toDomain() domain.DockerTaskDefinition {
	return domain.DockerTaskDefinition{
		ImageName:  d.ImageName,
		ContextDir: d.ContextDir,
	}
}

func fromDomainDockerTaskDef(def domain.DockerTaskDefinition) (dockerTaskDefDTO, error) {
	return dockerTaskDefDTO{
		ImageName:  def.ImageName,
		ContextDir: def.ContextDir,
	}, nil
}

type commandTaskDefDTO struct {
	Command    string `json:"command"`
	WorkingDir string `json:"working_dir"`
}

func (c commandTaskDefDTO) toDomain() domain.CommandTaskDefinition {
	return domain.CommandTaskDefinition{
		Command:    c.Command,
		WorkingDir: c.WorkingDir,
	}
}

func fromDomainCommandTaskDef(def domain.CommandTaskDefinition) (commandTaskDefDTO, error) {
	return commandTaskDefDTO{
		Command:    def.Command,
		WorkingDir: def.WorkingDir,
	}, nil
}

type jsonTaskRepository struct {
	config   JsonRepository
	fileName string
}

func NewJSONTaskRepository(repo *JsonRepository) domain.TaskRepository {
	path := filepath.Join(repo.RootDir, "tasks.json")
	if _, err := os.Stat(path); os.IsNotExist(err) {
		os.WriteFile(path, []byte("[]"), 0644)
	}

	return &jsonTaskRepository{
		config:   *repo,
		fileName: path,
	}
}

func (r *jsonTaskRepository) readCurrent() ([]domain.Task, error) {
	f, err := os.Open(r.fileName)
	if err != nil {
		return nil, err
	}

	var tasks []taskDTO
	err = json.NewDecoder(f).Decode(&tasks)
	if err != nil {
		return nil, err
	}
	err = f.Close()
	if err != nil {
		return nil, err
	}

	domainTasks := make([]domain.Task, len(tasks))
	for i, task := range tasks {
		domainTask, err := task.toDomain()
		if err != nil {
			return nil, err
		}
		domainTasks[i] = domainTask
	}

	return domainTasks, nil
}

func (r *jsonTaskRepository) write(tasks []domain.Task) error {
	dtoTasks := make([]taskDTO, 0, len(tasks))
	for _, task := range tasks {
		dtoTask, err := fromDomainTask(task)
		if err != nil {
			return err
		}
		dtoTasks = append(dtoTasks, dtoTask)
	}

	data, err := json.MarshalIndent(dtoTasks, "", "  ")

	if err != nil {
		return err
	}

	tmpFileName := r.fileName + ".tmp"
	err = os.WriteFile(tmpFileName, data, 0644)
	if err != nil {
		return err
	}

	err = os.Rename(tmpFileName, r.fileName)
	if err != nil {
		return err
	}

	return nil
}

func (r *jsonTaskRepository) CreateTask(task domain.Task) (domain.Task, error) {
	tasks, err := r.readCurrent()
	if err != nil {
		return domain.Task{}, err
	}

	if slices.ContainsFunc(tasks, func(t domain.Task) bool {
		return t.ID == task.ID
	}) {
		return domain.Task{}, errors.New("task with the same ID already exists")
	}
	if slices.ContainsFunc(tasks, func(t domain.Task) bool {
		return t.Name == task.Name
	}) {
		return domain.Task{}, errors.New("task with the same name already exists")
	}

	tasks = append(tasks, task)

	err = r.write(tasks)
	if err != nil {
		return domain.Task{}, err
	}

	return task, nil
}

func (r *jsonTaskRepository) GetTaskByID(id string) (domain.Task, error) {
	tasks, err := r.readCurrent()
	if err != nil {
		return domain.Task{}, err
	}

	taskIndex := slices.IndexFunc(tasks, func(t domain.Task) bool {
		return t.ID == id
	})
	if taskIndex == -1 {
		return domain.Task{}, domain.ErrNotFound
	}

	return tasks[taskIndex], nil
}

func (r *jsonTaskRepository) GetTaskByName(name string) (domain.Task, error) {
	tasks, err := r.readCurrent()
	if err != nil {
		return domain.Task{}, err
	}

	taskIndex := slices.IndexFunc(tasks, func(t domain.Task) bool {
		return t.Name == name
	})
	if taskIndex == -1 {
		return domain.Task{}, domain.ErrNotFound
	}

	return tasks[taskIndex], nil
}

func (r *jsonTaskRepository) GetAllTasks() ([]domain.Task, error) {
	tasks, err := r.readCurrent()
	if err != nil {
		return nil, err
	}
	return tasks, nil
}

func (r *jsonTaskRepository) UpdateTask(task domain.Task) (domain.Task, error) {
	tasks, err := r.readCurrent()
	if err != nil {
		return domain.Task{}, err
	}

	taskIndex := slices.IndexFunc(tasks, func(t domain.Task) bool {
		return t.ID == task.ID
	})
	if taskIndex == -1 {
		return domain.Task{}, domain.ErrNotFound
	}

	taskWithSameNameIndex := slices.IndexFunc(tasks, func(t domain.Task) bool {
		return t.Name == task.Name && t.ID != task.ID
	})
	if taskWithSameNameIndex != -1 {
		return domain.Task{}, errors.New("task with the same name already exists")
	}

	tasks[taskIndex] = task

	err = r.write(tasks)
	if err != nil {
		return domain.Task{}, err
	}

	return task, nil
}

func (r *jsonTaskRepository) DeleteTask(id string) error {
	tasks, err := r.readCurrent()
	if err != nil {
		return err
	}

	newTasks := slices.DeleteFunc(tasks, func(t domain.Task) bool {
		return t.ID == id
	})
	if len(newTasks) == len(tasks) {
		return domain.ErrNotFound
	}

	err = r.write(newTasks)
	if err != nil {
		return err
	}

	return nil
}
