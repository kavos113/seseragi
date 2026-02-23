package json

import (
	"encoding/json"
	"os"
	"path/filepath"
	"slices"

	"github.com/kavos113/seseragi/model"
)

type jsonWorkflowRunRepository struct {
	config   JsonRepository
	fileName string
}

func NewJSONWorkflowRunRepository(repo *JsonRepository) model.WorkflowRunRepository {
	path := filepath.Join(repo.RootDir, "workflow_runs.json")
	if _, err := os.Stat(path); os.IsNotExist(err) {
		os.WriteFile(path, []byte("[]"), 0644)
	}	

	return &jsonWorkflowRunRepository{
		config:   *repo,
		fileName: path,
	}
}

func (r *jsonWorkflowRunRepository) CreateWorkflowRun(workflowRun model.WorkflowRun) (model.WorkflowRun, error) {
	f, err := os.OpenFile(r.fileName, os.O_RDWR|os.O_CREATE, 0755)
	if err != nil {
		return model.WorkflowRun{}, err
	}

	var workflowRuns []model.WorkflowRun
	err = json.NewDecoder(f).Decode(&workflowRuns)
	if err != nil {
		return model.WorkflowRun{}, err
	}
	err = f.Close()
	if err != nil {
		return model.WorkflowRun{}, err
	}

	workflowRuns = append(workflowRuns, workflowRun)

	data, err := json.Marshal(workflowRuns)
	if err != nil {
		return model.WorkflowRun{}, err
	}

	tmpFileName := r.fileName + ".tmp"
	err = os.WriteFile(tmpFileName, data, 0644)
	if err != nil {
		return model.WorkflowRun{}, err
	}

	err = os.Rename(tmpFileName, r.fileName)
	if err != nil {
		return model.WorkflowRun{}, err
	}

	return workflowRun, nil
}

func (r *jsonWorkflowRunRepository) GetWorkflowRunByID(id string) (model.WorkflowRun, error) {
	f, err := os.Open(r.fileName)
	if err != nil {
		return model.WorkflowRun{}, err
	}

	var workflowRuns []model.WorkflowRun
	err = json.NewDecoder(f).Decode(&workflowRuns)
	if err != nil {
		return model.WorkflowRun{}, err
	}
	err = f.Close()
	if err != nil {
		return model.WorkflowRun{}, err
	}

	index := slices.IndexFunc(workflowRuns, func(wr model.WorkflowRun) bool {
		return wr.ID == id
	})
	if index == -1 {
		return model.WorkflowRun{}, model.ErrNotFound
	}

	return workflowRuns[index], nil
}