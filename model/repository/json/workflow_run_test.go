package json

import (
	"testing"
	"time"

	"github.com/kavos113/seseragi/model"
	"github.com/stretchr/testify/assert"
)

func TestCreateWorkflowRun(t *testing.T) {
	tests := []struct {
		name        string
		workflowRun model.WorkflowRun
		expected    model.WorkflowRun
		wantErr     bool
	}{
		{
			name: "success: create workflow run",
			workflowRun: model.WorkflowRun{
				ID:         "1",
				WorkflowID: "workflow-1",
				StartTime:  time.Date(2024, 6, 1, 12, 0, 0, 0, time.UTC),
				EndTime:    time.Date(2024, 6, 1, 12, 30, 0, 0, time.UTC),
				Status:     model.WorkflowStatusCompleted,
			},
			expected: model.WorkflowRun{
				ID:         "1",
				WorkflowID: "workflow-1",
				StartTime:  time.Date(2024, 6, 1, 12, 0, 0, 0, time.UTC),
				EndTime:    time.Date(2024, 6, 1, 12, 30, 0, 0, time.UTC),
				Status:     model.WorkflowStatusCompleted,
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := setupTestRepository(t)
			wf := NewJSONWorkflowRunRepository(repo)

			created, err := wf.CreateWorkflowRun(tt.workflowRun)
			if err != nil {
				if !tt.wantErr {
					t.Errorf("CreateWorkflowRun() error = %v, wantErr %v", err, tt.wantErr)
				}
				return
			}

			assert.Equal(t, tt.expected, created)
		})
	}
}

func TestGetAllWorkflowRuns(t *testing.T) {
	tests := []struct {
		name         string
		existingRuns []model.WorkflowRun
		expected     []model.WorkflowRun
		wantErr      bool
	}{
		{
			name: "success: get all workflow runs",
			existingRuns: []model.WorkflowRun{
				{
					ID:         "1",
					WorkflowID: "workflow-1",
					StartTime:  time.Date(2024, 6, 1, 12, 0, 0, 0, time.UTC),
					EndTime:    time.Date(2024, 6, 1, 12, 30, 0, 0, time.UTC),
					Status:     model.WorkflowStatusCompleted,
				},
				{
					ID:         "2",
					WorkflowID: "workflow-2",
					StartTime:  time.Date(2024, 6, 2, 14, 0, 0, 0, time.UTC),
					EndTime:    time.Date(2024, 6, 2, 14, 45, 0, 0, time.UTC),
					Status:     model.WorkflowStatusFailed,
				},
			},
			expected: []model.WorkflowRun{
				{
					ID:         "1",
					WorkflowID: "workflow-1",
					StartTime:  time.Date(2024, 6, 1, 12, 0, 0, 0, time.UTC),
					EndTime:    time.Date(2024, 6, 1, 12, 30, 0, 0, time.UTC),
					Status:     model.WorkflowStatusCompleted,
				},
				{
					ID:         "2",
					WorkflowID: "workflow-2",
					StartTime:  time.Date(2024, 6, 2, 14, 0, 0, 0, time.UTC),
					EndTime:    time.Date(2024, 6, 2, 14, 45, 0, 0, time.UTC),
					Status:     model.WorkflowStatusFailed,
				},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := setupTestRepository(t)
			wf := NewJSONWorkflowRunRepository(repo)

			for _, run := range tt.existingRuns {
				_, err := wf.CreateWorkflowRun(run)
				if err != nil {
					t.Fatalf("Failed to create workflow run: %v", err)
				}
			}

			all, err := wf.GetAllWorkflowRuns()
			if err != nil {
				if !tt.wantErr {
					t.Errorf("GetAllWorkflowRuns() error = %v, wantErr %v", err, tt.wantErr)
				}
				return
			}

			assert.Equal(t, tt.expected, all)
		})
	}
}

func TestGetWorkflowRunByID(t *testing.T) {
	tests := []struct {
		name         string
		existingRuns []model.WorkflowRun
		id           string
		expected     model.WorkflowRun
		wantErr      bool
	}{
		{
			name: "success: get workflow run by ID",
			existingRuns: []model.WorkflowRun{
				{
					ID:         "1",
					WorkflowID: "workflow-1",
					StartTime:  time.Date(2024, 6, 1, 12, 0, 0, 0, time.UTC),
					EndTime:    time.Date(2024, 6, 1, 12, 30, 0, 0, time.UTC),
					Status:     model.WorkflowStatusCompleted,
				},
			},
			id: "1",
			expected: model.WorkflowRun{
				ID:         "1",
				WorkflowID: "workflow-1",
				StartTime:  time.Date(2024, 6, 1, 12, 0, 0, 0, time.UTC),
				EndTime:    time.Date(2024, 6, 1, 12, 30, 0, 0, time.UTC),
				Status:     model.WorkflowStatusCompleted,
			},
			wantErr: false,
		},
		{
			name: "failure: workflow run not found",
			existingRuns: []model.WorkflowRun{
				{
					ID:         "1",
					WorkflowID: "workflow-1",
					StartTime:  time.Date(2024, 6, 1, 12, 0, 0, 0, time.UTC),
					EndTime:    time.Date(2024, 6, 1, 12, 30, 0, 0, time.UTC),
					Status:     model.WorkflowStatusCompleted,
				},
			},
			id:       "non-existent-id",
			expected: model.WorkflowRun{},
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := setupTestRepository(t)
			wf := NewJSONWorkflowRunRepository(repo)

			for _, run := range tt.existingRuns {
				_, err := wf.CreateWorkflowRun(run)
				if err != nil {
					t.Fatalf("Failed to create workflow run: %v", err)
				}
			}

			got, err := wf.GetWorkflowRunByID(tt.id)
			if err != nil {
				if !tt.wantErr {
					t.Errorf("GetWorkflowRunByID() error = %v, wantErr %v", err, tt.wantErr)
				}
				return
			}

			assert.Equal(t, tt.expected, got)
		})
	}
}

func TestGetWorkflowRunsByWorkflowID(t *testing.T) {
	tests := []struct {
		name         string
		existingRuns []model.WorkflowRun
		workflowID   string
		expected     []model.WorkflowRun
		wantErr      bool
	}{
		{
			name: "success: get workflow runs by workflow ID",
			existingRuns: []model.WorkflowRun{
				{
					ID:         "1",
					WorkflowID: "workflow-1",
					StartTime:  time.Date(2024, 6, 1, 12, 0, 0, 0, time.UTC),
					EndTime:    time.Date(2024, 6, 1, 12, 30, 0, 0, time.UTC),
					Status:     model.WorkflowStatusCompleted,
				},
				{
					ID:         "2",
					WorkflowID: "workflow-1",
					StartTime:  time.Date(2024, 6, 2, 14, 0, 0, 0, time.UTC),
					EndTime:    time.Date(2024, 6, 2, 14, 45, 0, 0, time.UTC),
					Status:     model.WorkflowStatusFailed,
				},
				{
					ID:         "3",
					WorkflowID: "workflow-2",
					StartTime:  time.Date(2024, 6, 3, 10, 0, 0, 0, time.UTC),
					EndTime:    time.Date(2024, 6, 3, 10, 30, 0, 0, time.UTC),
					Status:     model.WorkflowStatusCompleted,
				},
			},
			workflowID: "workflow-1",
			expected: []model.WorkflowRun{
				{
					ID:         "1",
					WorkflowID: "workflow-1",
					StartTime:  time.Date(2024, 6, 1, 12, 0, 0, 0, time.UTC),
					EndTime:    time.Date(2024, 6, 1, 12, 30, 0, 0, time.UTC),
					Status:     model.WorkflowStatusCompleted,
				},
				{
					ID:         "2",
					WorkflowID: "workflow-1",
					StartTime:  time.Date(2024, 6, 2, 14, 0, 0, 0, time.UTC),
					EndTime:    time.Date(2024, 6, 2, 14, 45, 0, 0, time.UTC),
					Status:     model.WorkflowStatusFailed,
				},
			},
			wantErr: false,
		},
		{
			name: "failure: no workflow runs for given workflow ID",
			existingRuns: []model.WorkflowRun{
				{
					ID:         "1",
					WorkflowID: "workflow-1",
					StartTime:  time.Date(2024, 6, 1, 12, 0, 0, 0, time.UTC),
					EndTime:    time.Date(2024, 6, 1, 12, 30, 0, 0, time.UTC),
					Status:     model.WorkflowStatusCompleted,
				},
			},
			workflowID: "non-existent-workflow-id",
			expected:   []model.WorkflowRun{},
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := setupTestRepository(t)
			wf := NewJSONWorkflowRunRepository(repo)

			for _, run := range tt.existingRuns {
				_, err := wf.CreateWorkflowRun(run)
				if err != nil {
					t.Fatalf("Failed to create workflow run: %v", err)
				}
			}

			got, err := wf.GetWorkflowRunsByWorkflowID(tt.workflowID)
			if err != nil {
				if !tt.wantErr {
					t.Errorf("GetWorkflowRunsByWorkflowID() error = %v, wantErr %v", err, tt.wantErr)
				}
				return
			}

			assert.Equal(t, tt.expected, got)
		})
	}
}

func TestGetWorkflowRunsAfter(t *testing.T) {
	tests := []struct {
		name         string
		existingRuns []model.WorkflowRun
		workflowID   string
		after        time.Time
		expected     []model.WorkflowRun
		wantErr      bool
	}{
		{
			name: "success: get workflow runs after a specific time",
			existingRuns: []model.WorkflowRun{
				{
					ID:         "1",
					WorkflowID: "workflow-1",
					StartTime:  time.Date(2024, 6, 1, 12, 0, 0, 0, time.UTC),
					EndTime:    time.Date(2024, 6, 1, 12, 30, 0, 0, time.UTC),
					Status:     model.WorkflowStatusCompleted,
				},
				{
					ID:         "2",
					WorkflowID: "workflow-1",
					StartTime:  time.Date(2024, 6, 2, 14, 0, 0, 0, time.UTC),
					EndTime:    time.Date(2024, 6, 2, 14, 45, 0, 0, time.UTC),
					Status:     model.WorkflowStatusFailed,
				},
				{
					ID:         "3",
					WorkflowID: "workflow-1",
					StartTime:  time.Date(2024, 6, 3, 10, 0, 0, 0, time.UTC),
					EndTime:    time.Date(2024, 6, 3, 10, 30, 0, 0, time.UTC),
					Status:     model.WorkflowStatusCompleted,
				},
			},
			workflowID: "workflow-1",
			after:      time.Date(2024, 6, 2, 0, 0, 0, 0, time.UTC),
			expected: []model.WorkflowRun{
				{
					ID:         "2",
					WorkflowID: "workflow-1",
					StartTime:  time.Date(2024, 6, 2, 14, 0, 0, 0, time.UTC),
					EndTime:    time.Date(2024, 6, 2, 14, 45, 0, 0, time.UTC),
					Status:     model.WorkflowStatusFailed,
				},
				{
					ID:         "3",
					WorkflowID: "workflow-1",
					StartTime:  time.Date(2024, 6, 3, 10, 0, 0, 0, time.UTC),
					EndTime:    time.Date(2024, 6, 3, 10, 30, 0, 0, time.UTC),
					Status:     model.WorkflowStatusCompleted,
				},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := setupTestRepository(t)
			wf := NewJSONWorkflowRunRepository(repo)

			for _, run := range tt.existingRuns {
				_, err := wf.CreateWorkflowRun(run)
				if err != nil {
					t.Fatalf("Failed to create workflow run: %v", err)
				}
			}

			got, err := wf.GetWorkflowRunsAfter(tt.workflowID, tt.after)
			if err != nil {
				if !tt.wantErr {
					t.Errorf("GetWorkflowRunsAfter() error = %v, wantErr %v", err, tt.wantErr)
				}
				return
			}

			assert.Equal(t, tt.expected, got)
		})
	}
}
