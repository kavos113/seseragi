package json

import (
	"encoding/json"
	"os"
	"path/filepath"
	"slices"

	"github.com/kavos113/seseragi/internal/domain"
)

type jsonWorkflowRepository struct {
	config   JsonRepository
	fileName string
}

func NewJSONWorkflowRepository(repo *JsonRepository) domain.WorkflowRepository {
	path := filepath.Join(repo.RootDir, "workflows.json")
	if _, err := os.Stat(path); os.IsNotExist(err) {
		os.WriteFile(path, []byte("[]"), 0644)
	}

	return &jsonWorkflowRepository{
		config:   *repo,
		fileName: path,
	}
}

func (r *jsonWorkflowRepository) readCurrent() ([]domain.Workflow, error) {
	f, err := os.Open(r.fileName)
	if err != nil {
		return nil, err
	}

	var workflows []domain.Workflow
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

func (r *jsonWorkflowRepository) write(workflows []domain.Workflow) error {
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

func (r *jsonWorkflowRepository) CreateWorkflow(workflow domain.Workflow) (domain.Workflow, error) {
	workflows, err := r.readCurrent()
	if err != nil {
		return domain.Workflow{}, err
	}

	workflows = append(workflows, workflow)

	err = r.write(workflows)
	if err != nil {
		return domain.Workflow{}, err
	}

	return workflow, nil
}

func (r *jsonWorkflowRepository) GetAllWorkflows() ([]domain.Workflow, error) {
	return r.readCurrent()
}

func (r *jsonWorkflowRepository) GetWorkflowByID(id string) (domain.Workflow, error) {
	workflows, err := r.readCurrent()
	if err != nil {
		return domain.Workflow{}, err
	}

	workflowIndex := slices.IndexFunc(workflows, func(w domain.Workflow) bool {
		return w.ID == id
	})
	if workflowIndex == -1 {
		return domain.Workflow{}, domain.ErrNotFound
	}
	return workflows[workflowIndex], nil
}

func (r *jsonWorkflowRepository) UpdateWorkflow(workflow domain.Workflow) (domain.Workflow, error) {
	workflows, err := r.readCurrent()
	if err != nil {
		return domain.Workflow{}, err
	}

	workflowIndex := slices.IndexFunc(workflows, func(w domain.Workflow) bool {
		return w.ID == workflow.ID
	})
	if workflowIndex == -1 {
		return domain.Workflow{}, domain.ErrNotFound
	}

	workflows[workflowIndex] = workflow

	err = r.write(workflows)
	if err != nil {
		return domain.Workflow{}, err
	}

	return workflow, nil
}

func (r *jsonWorkflowRepository) AddNodeToWorkflow(workflowID string, node domain.Node) (domain.Workflow, error) {
	workflows, err := r.readCurrent()
	if err != nil {
		return domain.Workflow{}, err
	}

	workflowIndex := slices.IndexFunc(workflows, func(w domain.Workflow) bool {
		return w.ID == workflowID
	})
	if workflowIndex == -1 {
		return domain.Workflow{}, domain.ErrNotFound
	}

	workflows[workflowIndex].Nodes = append(workflows[workflowIndex].Nodes, node)

	err = r.write(workflows)
	if err != nil {
		return domain.Workflow{}, err
	}

	return workflows[workflowIndex], nil
}

func (r *jsonWorkflowRepository) DeleteNodeFromWorkflow(workflowID string, nodeName string) (domain.Workflow, error) {
	workflows, err := r.readCurrent()
	if err != nil {
		return domain.Workflow{}, err
	}

	workflowIndex := slices.IndexFunc(workflows, func(w domain.Workflow) bool {
		return w.ID == workflowID
	})
	if workflowIndex == -1 {
		return domain.Workflow{}, domain.ErrNotFound
	}
	workflow := &workflows[workflowIndex]

	newNodes := slices.DeleteFunc(workflow.Nodes, func(n domain.Node) bool {
		return n.Name == nodeName
	})
	if len(newNodes) == len(workflow.Nodes) {
		return domain.Workflow{}, domain.ErrNotFound
	}
	workflow.Nodes = newNodes

	err = r.write(workflows)
	if err != nil {
		return domain.Workflow{}, err
	}

	return *workflow, nil
}

func (r *jsonWorkflowRepository) DeleteWorkflow(id string) error {
	workflows, err := r.readCurrent()
	if err != nil {
		return err
	}

	newWorkflows := slices.DeleteFunc(workflows, func(w domain.Workflow) bool {
		return w.ID == id
	})
	if len(newWorkflows) == len(workflows) {
		return domain.ErrNotFound
	}

	err = r.write(newWorkflows)
	if err != nil {
		return err
	}

	return nil
}
