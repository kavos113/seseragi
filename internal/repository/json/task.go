package json

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"slices"

	"github.com/kavos113/seseragi/internal/domain"
)

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

	var tasks []domain.Task
	err = json.NewDecoder(f).Decode(&tasks)
	if err != nil {
		return nil, err
	}
	err = f.Close()
	if err != nil {
		return nil, err
	}

	return tasks, nil
}

func (r *jsonTaskRepository) write(tasks []domain.Task) error {
	data, err := json.Marshal(tasks)
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
