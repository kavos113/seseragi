package json

import (
	"encoding/json"
	"os"
	"path/filepath"
	"slices"

	"github.com/kavos113/seseragi/model"
)

type jsonWorkflowRepository struct {
	config   JsonRepository
	fileName string
}

func NewJSONWorkflowRepository(repo *JsonRepository) model.WorkflowRepository {
	path := filepath.Join(repo.RootDir, "workflows.json")
	if _, err := os.Stat(path); os.IsNotExist(err) {
		os.WriteFile(path, []byte("[]"), 0644)
	}

	return &jsonWorkflowRepository{
		config:   *repo,
		fileName: path,
	}
}

func (r *jsonWorkflowRepository) readCurrent() ([]model.Workflow, error) {
	f, err := os.Open(r.fileName)
	if err != nil {
		return nil, err
	}

	var workflows []model.Workflow
	err = json.NewDecoder(f).Decode(&workflows)
	if err != nil {
		return nil, err
	}

	err = f.Close()
	if err != nil {
		return nil, err
	}

	return workflows, nil
}

func (r *jsonWorkflowRepository) write(workflows []model.Workflow) error {
	data, err := json.Marshal(workflows)
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

func (r *jsonWorkflowRepository) CreateWorkflow(workflow model.Workflow) (model.Workflow, error) {
	workflows, err := r.readCurrent()
	if err != nil {
		return model.Workflow{}, err
	}

	workflows = append(workflows, workflow)

	err = r.write(workflows)
	if err != nil {
		return model.Workflow{}, err
	}

	return workflow, nil
}

func (r *jsonWorkflowRepository) GetAllWorkflows() ([]model.Workflow, error) {
	return r.readCurrent()
}

func (r *jsonWorkflowRepository) GetWorkflowByID(id string) (model.Workflow, error) {
	workflows, err := r.readCurrent()
	if err != nil {
		return model.Workflow{}, err
	}

	workflowIndex := slices.IndexFunc(workflows, func(w model.Workflow) bool {
		return w.ID == id
	})
	if workflowIndex == -1 {
		return model.Workflow{}, model.ErrNotFound
	}
	return workflows[workflowIndex], nil
}

func (r *jsonWorkflowRepository) UpdateWorkflow(workflow model.Workflow) (model.Workflow, error) {
	workflows, err := r.readCurrent()
	if err != nil {
		return model.Workflow{}, err
	}

	workflowIndex := slices.IndexFunc(workflows, func(w model.Workflow) bool {
		return w.ID == workflow.ID
	})
	if workflowIndex == -1 {
		return model.Workflow{}, model.ErrNotFound
	}

	workflows[workflowIndex] = workflow

	err = r.write(workflows)
	if err != nil {
		return model.Workflow{}, err
	}

	return workflow, nil
}

func (r *jsonWorkflowRepository) AddNodeToWorkflow(workflowID string, node model.Node) (model.Workflow, error) {
	workflows, err := r.readCurrent()
	if err != nil {
		return model.Workflow{}, err
	}

	workflowIndex := slices.IndexFunc(workflows, func(w model.Workflow) bool {
		return w.ID == workflowID
	})
	if workflowIndex == -1 {
		return model.Workflow{}, model.ErrNotFound
	}

	workflows[workflowIndex].Nodes = append(workflows[workflowIndex].Nodes, node)

	err = r.write(workflows)
	if err != nil {
		return model.Workflow{}, err
	}

	return workflows[workflowIndex], nil
}

func (r *jsonWorkflowRepository) DeleteNodeFromWorkflow(workflowID string, taskID string) (model.Workflow, error) {
	workflows, err := r.readCurrent()
	if err != nil {
		return model.Workflow{}, err
	}

	workflowIndex := slices.IndexFunc(workflows, func(w model.Workflow) bool {
		return w.ID == workflowID
	})
	if workflowIndex == -1 {
		return model.Workflow{}, model.ErrNotFound
	}
	workflow := workflows[workflowIndex]

	newNodes := slices.DeleteFunc(workflow.Nodes, func(n model.Node) bool {
		return n.TaskID == taskID
	})
	if len(newNodes) == len(workflow.Nodes) {
		return model.Workflow{}, model.ErrNotFound
	}
	workflow.Nodes = newNodes

	err = r.write(workflows)
	if err != nil {
		return model.Workflow{}, err
	}

	return workflow, nil
}

func (r *jsonWorkflowRepository) DeleteWorkflow(id string) error {
	workflows, err := r.readCurrent()
	if err != nil {
		return err
	}

	newWorkflows := slices.DeleteFunc(workflows, func(w model.Workflow) bool {
		return w.ID == id
	})
	if len(newWorkflows) == len(workflows) {
		return model.ErrNotFound
	}

	err = r.write(newWorkflows)
	if err != nil {
		return err
	}

	return nil
}
