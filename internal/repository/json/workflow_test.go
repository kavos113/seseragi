package json

import (
	"testing"

	"github.com/kavos113/seseragi/internal/domain"
	"github.com/stretchr/testify/assert"
)

func TestCreateWorkflow(t *testing.T) {
	tests := []struct {
		name     string
		workflow domain.Workflow
		expected domain.Workflow
		wantErr  bool
	}{
		{
			name: "success",
			workflow: domain.Workflow{
				ID:   "1",
				Name: "Test Workflow",
			},
			expected: domain.Workflow{
				ID:   "1",
				Name: "Test Workflow",
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := setupTestRepository(t)
			wr := NewJSONWorkflowRepository(repo)

			result, err := wr.CreateWorkflow(tt.workflow)
			if (err != nil) != tt.wantErr {
				t.Errorf("CreateWorkflow() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestGetAllWorkflows(t *testing.T) {
	tests := []struct {
		name      string
		workflows []domain.Workflow
		expected  []domain.Workflow
		wantErr   bool
	}{
		{
			name: "success",
			workflows: []domain.Workflow{
				{ID: "1", Name: "Workflow 1"},
				{ID: "2", Name: "Workflow 2"},
			},
			expected: []domain.Workflow{
				{ID: "1", Name: "Workflow 1"},
				{ID: "2", Name: "Workflow 2"},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := setupTestRepository(t)
			wr := NewJSONWorkflowRepository(repo)

			for _, w := range tt.workflows {
				_, err := wr.CreateWorkflow(w)
				if err != nil {
					t.Fatalf("Failed to create workflow: %v", err)
				}
			}

			result, err := wr.GetAllWorkflows()
			if (err != nil) != tt.wantErr {
				t.Errorf("GetAllWorkflows() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestGetWorkflowByID(t *testing.T) {
	tests := []struct {
		name      string
		workflows []domain.Workflow
		id        string
		expected  domain.Workflow
		wantErr   bool
	}{
		{
			name: "success",
			workflows: []domain.Workflow{
				{ID: "1", Name: "Workflow 1"},
				{ID: "2", Name: "Workflow 2"},
			},
			id:       "1",
			expected: domain.Workflow{ID: "1", Name: "Workflow 1"},
			wantErr:  false,
		},
		{
			name: "failure: not found",
			workflows: []domain.Workflow{
				{ID: "1", Name: "Workflow 1"},
			},
			id:       "2",
			expected: domain.Workflow{},
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := setupTestRepository(t)
			wr := NewJSONWorkflowRepository(repo)

			for _, w := range tt.workflows {
				_, err := wr.CreateWorkflow(w)
				if err != nil {
					t.Fatalf("Failed to create workflow: %v", err)
				}
			}

			result, err := wr.GetWorkflowByID(tt.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetWorkflowByID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestUpdateWorkflow(t *testing.T) {
	tests := []struct {
		name      string
		workflows []domain.Workflow
		update    domain.Workflow
		expected  domain.Workflow
		wantErr   bool
	}{
		{
			name: "success",
			workflows: []domain.Workflow{
				{ID: "1", Name: "Workflow 1"},
			},
			update:   domain.Workflow{ID: "1", Name: "Updated Workflow 1"},
			expected: domain.Workflow{ID: "1", Name: "Updated Workflow 1"},
			wantErr:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := setupTestRepository(t)
			wr := NewJSONWorkflowRepository(repo)

			for _, w := range tt.workflows {
				_, err := wr.CreateWorkflow(w)
				if err != nil {
					t.Fatalf("Failed to create workflow: %v", err)
				}
			}

			result, err := wr.UpdateWorkflow(tt.update)
			if (err != nil) != tt.wantErr {
				t.Errorf("UpdateWorkflow() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			assert.Equal(t, tt.expected, result)

			// Verify the update is persisted
			persisted, err := wr.GetWorkflowByID(tt.update.ID)
			if err != nil {
				t.Fatalf("Failed to get workflow by ID: %v", err)
			}
			assert.Equal(t, tt.expected, persisted)
		})
	}
}

func TestAddNodeToWorkflow(t *testing.T) {
	tests := []struct {
		name       string
		workflows  []domain.Workflow
		workflowID string
		node       domain.Node
		expected   domain.Workflow
		wantErr    bool
	}{
		{
			name: "success",
			workflows: []domain.Workflow{
				{ID: "1", Name: "Workflow 1", Nodes: []domain.Node{}},
			},
			workflowID: "1",
			node:       domain.Node{Name: "Node 1"},
			expected: domain.Workflow{
				ID:   "1",
				Name: "Workflow 1",
				Nodes: []domain.Node{
					{Name: "Node 1"},
				},
			},
			wantErr: false,
		},
		{
			name: "failure: workflow not found",
			workflows: []domain.Workflow{
				{ID: "1", Name: "Workflow 1", Nodes: []domain.Node{}},
			},
			workflowID: "2",
			node:       domain.Node{Name: "Node 1"},
			expected:   domain.Workflow{},
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := setupTestRepository(t)
			wr := NewJSONWorkflowRepository(repo)

			for _, w := range tt.workflows {
				_, err := wr.CreateWorkflow(w)
				if err != nil {
					t.Fatalf("Failed to create workflow: %v", err)
				}
			}

			result, err := wr.AddNodeToWorkflow(tt.workflowID, tt.node)
			if (err != nil) != tt.wantErr {
				t.Errorf("AddNodeToWorkflow() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			assert.Equal(t, tt.expected, result)

			if !tt.wantErr {
				// Verify the node is added to the workflow
				persisted, err := wr.GetWorkflowByID(tt.workflowID)
				if err != nil {
					t.Fatalf("Failed to get workflow by ID: %v", err)
				}
				assert.Equal(t, tt.expected, persisted)
			}
		})
	}
}

func TestDeleteNodeFromWorkflow(t *testing.T) {
	tests := []struct {
		name       string
		workflows  []domain.Workflow
		workflowID string
		nodeID     string
		expected   domain.Workflow
		wantErr    bool
	}{
		{
			name: "success",
			workflows: []domain.Workflow{
				{
					ID:    "1",
					Name:  "Workflow 1",
					Nodes: []domain.Node{{Name: "Node 1"}}},
			},
			workflowID: "1",
			nodeID:     "Node 1",
			expected: domain.Workflow{
				ID:    "1",
				Name:  "Workflow 1",
				Nodes: []domain.Node{},
			},
			wantErr: false,
		},
		{
			name: "failure: workflow not found",
			workflows: []domain.Workflow{
				{ID: "1", Name: "Workflow 1", Nodes: []domain.Node{{Name: "Node 1"}}},
			},
			workflowID: "2",
			nodeID:     "Node 1",
			expected:   domain.Workflow{},
			wantErr:    true,
		},
		{
			name: "failure: node not found",
			workflows: []domain.Workflow{
				{ID: "1", Name: "Workflow 1", Nodes: []domain.Node{{Name: "Node 1"}}},
			},
			workflowID: "1",
			nodeID:     "Node 2",
			expected:   domain.Workflow{},
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := setupTestRepository(t)
			wr := NewJSONWorkflowRepository(repo)

			for _, w := range tt.workflows {
				_, err := wr.CreateWorkflow(w)
				if err != nil {
					t.Fatalf("Failed to create workflow: %v", err)
				}
			}

			result, err := wr.DeleteNodeFromWorkflow(tt.workflowID, tt.nodeID)
			if (err != nil) != tt.wantErr {
				t.Errorf("DeleteNodeFromWorkflow() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			assert.Equal(t, tt.expected, result)

			if !tt.wantErr {
				// Verify the node is deleted from the workflow
				persisted, err := wr.GetWorkflowByID(tt.workflowID)
				if err != nil {
					t.Fatalf("Failed to get workflow by ID: %v", err)
				}
				assert.Equal(t, tt.expected, persisted)
			}
		})
	}
}

func TestDeleteWorkflow(t *testing.T) {
	tests := []struct {
		name      string
		workflows []domain.Workflow
		id        string
		wantErr   bool
	}{
		{
			name: "success",
			workflows: []domain.Workflow{
				{ID: "1", Name: "Workflow 1"},
				{ID: "2", Name: "Workflow 2"},
			},
			id:      "1",
			wantErr: false,
		},
		{
			name: "failure: not found",
			workflows: []domain.Workflow{
				{ID: "1", Name: "Workflow 1"},
			},
			id:      "2",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := setupTestRepository(t)
			wr := NewJSONWorkflowRepository(repo)

			for _, w := range tt.workflows {
				_, err := wr.CreateWorkflow(w)
				if err != nil {
					t.Fatalf("Failed to create workflow: %v", err)
				}
			}

			err := wr.DeleteWorkflow(tt.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("DeleteWorkflow() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				// Verify the workflow is deleted
				_, err := wr.GetWorkflowByID(tt.id)
				assert.Error(t, err)
			}
		})
	}
}
