package json

import (
	"encoding/json"
	"os"
	"path/filepath"
	"slices"
	"time"

	"github.com/kavos113/seseragi/internal/domain"
)

type jsonWorkflowRunRepository struct {
	config   JsonRepository
	fileName string
}

func NewJSONWorkflowRunRepository(repo *JsonRepository) domain.WorkflowRunRepository {
	path := filepath.Join(repo.RootDir, "workflow_runs.json")
	if _, err := os.Stat(path); os.IsNotExist(err) {
		os.WriteFile(path, []byte("[]"), 0644)
	}

	return &jsonWorkflowRunRepository{
		config:   *repo,
		fileName: path,
	}
}

func (r *jsonWorkflowRunRepository) readCurrent() ([]domain.WorkflowRun, error) {
	f, err := os.Open(r.fileName)
	if err != nil {
		return nil, err
	}

	var workflowRuns []domain.WorkflowRun
	err = json.NewDecoder(f).Decode(&workflowRuns)
	if err != nil {
		return nil, err
	}
	err = f.Close()
	if err != nil {
		return nil, err
	}

	return workflowRuns, nil
}

func (r *jsonWorkflowRunRepository) write(workflowRuns []domain.WorkflowRun) error {
	data, err := json.Marshal(workflowRuns)
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

func (r *jsonWorkflowRunRepository) CreateWorkflowRun(workflowRun domain.WorkflowRun) (domain.WorkflowRun, error) {
	workflowRuns, err := r.readCurrent()
	if err != nil {
		return domain.WorkflowRun{}, err
	}

	workflowRuns = append(workflowRuns, workflowRun)

	err = r.write(workflowRuns)
	if err != nil {
		return domain.WorkflowRun{}, err
	}
	return workflowRun, nil
}

func (r *jsonWorkflowRunRepository) GetAllWorkflowRuns() ([]domain.WorkflowRun, error) {
	return r.readCurrent()
}

func (r *jsonWorkflowRunRepository) GetWorkflowRunByID(id string) (domain.WorkflowRun, error) {
	workflowRuns, err := r.readCurrent()
	if err != nil {
		return domain.WorkflowRun{}, err
	}

	index := slices.IndexFunc(workflowRuns, func(wr domain.WorkflowRun) bool {
		return wr.ID == id
	})
	if index == -1 {
		return domain.WorkflowRun{}, domain.ErrNotFound
	}

	return workflowRuns[index], nil
}

func (r *jsonWorkflowRunRepository) GetWorkflowRunsByWorkflowID(workflowID string) ([]domain.WorkflowRun, error) {
	workflowRuns, err := r.readCurrent()
	if err != nil {
		return nil, err
	}

	var result []domain.WorkflowRun
	for _, wr := range workflowRuns {
		if wr.WorkflowID == workflowID {
			result = append(result, wr)
		}
	}
	if len(result) == 0 {
		return nil, domain.ErrNotFound
	}

	return result, nil
}

func (r *jsonWorkflowRunRepository) GetWorkflowRunsAfter(workflowID string, after time.Time) ([]domain.WorkflowRun, error) {
	workflowRuns, err := r.readCurrent()
	if err != nil {
		return nil, err
	}

	var result []domain.WorkflowRun
	for _, wr := range workflowRuns {
		if wr.WorkflowID == workflowID && wr.StartTime.After(after) {
			result = append(result, wr)
		}
	}
	if len(result) == 0 {
		return nil, domain.ErrNotFound
	}

	return result, nil
}
