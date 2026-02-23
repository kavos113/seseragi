package json

import (
	"encoding/json"
	"os"
	"path/filepath"
	"slices"

	"github.com/kavos113/seseragi/model"
)

type jsonTaskRepository struct {
	config   JsonRepository
	fileName string
}

func NewJSONTaskRepository(repo *JsonRepository) model.TaskRepository {
	path := filepath.Join(repo.RootDir, "tasks.json")
	if _, err := os.Stat(path); os.IsNotExist(err) {
		os.WriteFile(path, []byte("[]"), 0644)
	}

	return &jsonTaskRepository{
		config:   *repo,
		fileName: path,
	}
}

func (r *jsonTaskRepository) readCurrent() ([]model.Task, error) {
	f, err := os.Open(r.fileName)
	if err != nil {
		return nil, err
	}

	var tasks []model.Task
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

func (r *jsonTaskRepository) write(tasks []model.Task) error {
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

func (r *jsonTaskRepository) CreateTask(task model.Task) (model.Task, error) {
	tasks, err := r.readCurrent()
	if err != nil {
		return model.Task{}, err
	}

	tasks = append(tasks, task)

	err = r.write(tasks)
	if err != nil {
		return model.Task{}, err
	}

	return task, nil
}

func (r *jsonTaskRepository) GetTaskByID(id string) (model.Task, error) {
	tasks, err := r.readCurrent()
	if err != nil {
		return model.Task{}, err
	}

	taskIndex := slices.IndexFunc(tasks, func(t model.Task) bool {
		return t.ID == id
	})
	if taskIndex == -1 {
		return model.Task{}, model.ErrNotFound
	}

	return tasks[taskIndex], nil
}

func (r *jsonTaskRepository) UpdateTask(task model.Task) (model.Task, error) {
	tasks, err := r.readCurrent()
	if err != nil {
		return model.Task{}, err
	}

	taskIndex := slices.IndexFunc(tasks, func(t model.Task) bool {
		return t.ID == task.ID
	})
	if taskIndex == -1 {
		return model.Task{}, model.ErrNotFound
	}

	tasks[taskIndex] = task

	err = r.write(tasks)
	if err != nil {
		return model.Task{}, err
	}

	return task, nil
}

func (r *jsonTaskRepository) DeleteTask(id string) error {
	tasks, err := r.readCurrent()
	if err != nil {
		return err
	}

	newTasks := slices.DeleteFunc(tasks, func(t model.Task) bool {
		return t.ID == id
	})
	if len(newTasks) == len(tasks) {
		return model.ErrNotFound
	}

	err = r.write(newTasks)
	if err != nil {
		return err
	}

	return nil
}
