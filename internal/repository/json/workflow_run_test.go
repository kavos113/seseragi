package json

import (
	"testing"
	"time"

	"github.com/kavos113/seseragi/internal/domain"
	"github.com/stretchr/testify/assert"
)

func TestCreateWorkflowRun(t *testing.T) {
	tests := []struct {
		name        string
		workflowRun domain.WorkflowRun
		expected    domain.WorkflowRun
		wantErr     bool
	}{
		{
			name: "success: create workflow run",
			workflowRun: domain.WorkflowRun{
				ID:         "1",
				WorkflowID: "workflow-1",
				StartTime:  time.Date(2024, 6, 1, 12, 0, 0, 0, time.UTC),
				EndTime:    time.Date(2024, 6, 1, 12, 30, 0, 0, time.UTC),
				Status:     domain.WorkflowStatusCompleted,
			},
			expected: domain.WorkflowRun{
				ID:         "1",
				WorkflowID: "workflow-1",
				StartTime:  time.Date(2024, 6, 1, 12, 0, 0, 0, time.UTC),
				EndTime:    time.Date(2024, 6, 1, 12, 30, 0, 0, time.UTC),
				Status:     domain.WorkflowStatusCompleted,
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
		existingRuns []domain.WorkflowRun
		expected     []domain.WorkflowRun
		wantErr      bool
	}{
		{
			name: "success: get all workflow runs",
			existingRuns: []domain.WorkflowRun{
				{
					ID:         "1",
					WorkflowID: "workflow-1",
					StartTime:  time.Date(2024, 6, 1, 12, 0, 0, 0, time.UTC),
					EndTime:    time.Date(2024, 6, 1, 12, 30, 0, 0, time.UTC),
					Status:     domain.WorkflowStatusCompleted,
				},
				{
					ID:         "2",
					WorkflowID: "workflow-2",
					StartTime:  time.Date(2024, 6, 2, 14, 0, 0, 0, time.UTC),
					EndTime:    time.Date(2024, 6, 2, 14, 45, 0, 0, time.UTC),
					Status:     domain.WorkflowStatusFailed,
				},
			},
			expected: []domain.WorkflowRun{
				{
					ID:         "1",
					WorkflowID: "workflow-1",
					StartTime:  time.Date(2024, 6, 1, 12, 0, 0, 0, time.UTC),
					EndTime:    time.Date(2024, 6, 1, 12, 30, 0, 0, time.UTC),
					Status:     domain.WorkflowStatusCompleted,
				},
				{
					ID:         "2",
					WorkflowID: "workflow-2",
					StartTime:  time.Date(2024, 6, 2, 14, 0, 0, 0, time.UTC),
					EndTime:    time.Date(2024, 6, 2, 14, 45, 0, 0, time.UTC),
					Status:     domain.WorkflowStatusFailed,
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
		existingRuns []domain.WorkflowRun
		id           string
		expected     domain.WorkflowRun
		wantErr      bool
	}{
		{
			name: "success: get workflow run by ID",
			existingRuns: []domain.WorkflowRun{
				{
					ID:         "1",
					WorkflowID: "workflow-1",
					StartTime:  time.Date(2024, 6, 1, 12, 0, 0, 0, time.UTC),
					EndTime:    time.Date(2024, 6, 1, 12, 30, 0, 0, time.UTC),
					Status:     domain.WorkflowStatusCompleted,
				},
			},
			id: "1",
			expected: domain.WorkflowRun{
				ID:         "1",
				WorkflowID: "workflow-1",
				StartTime:  time.Date(2024, 6, 1, 12, 0, 0, 0, time.UTC),
				EndTime:    time.Date(2024, 6, 1, 12, 30, 0, 0, time.UTC),
				Status:     domain.WorkflowStatusCompleted,
			},
			wantErr: false,
		},
		{
			name: "failure: workflow run not found",
			existingRuns: []domain.WorkflowRun{
				{
					ID:         "1",
					WorkflowID: "workflow-1",
					StartTime:  time.Date(2024, 6, 1, 12, 0, 0, 0, time.UTC),
					EndTime:    time.Date(2024, 6, 1, 12, 30, 0, 0, time.UTC),
					Status:     domain.WorkflowStatusCompleted,
				},
			},
			id:       "non-existent-id",
			expected: domain.WorkflowRun{},
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
		existingRuns []domain.WorkflowRun
		workflowID   string
		expected     []domain.WorkflowRun
		wantErr      bool
	}{
		{
			name: "success: get workflow runs by workflow ID",
			existingRuns: []domain.WorkflowRun{
				{
					ID:         "1",
					WorkflowID: "workflow-1",
					StartTime:  time.Date(2024, 6, 1, 12, 0, 0, 0, time.UTC),
					EndTime:    time.Date(2024, 6, 1, 12, 30, 0, 0, time.UTC),
					Status:     domain.WorkflowStatusCompleted,
				},
				{
					ID:         "2",
					WorkflowID: "workflow-1",
					StartTime:  time.Date(2024, 6, 2, 14, 0, 0, 0, time.UTC),
					EndTime:    time.Date(2024, 6, 2, 14, 45, 0, 0, time.UTC),
					Status:     domain.WorkflowStatusFailed,
				},
				{
					ID:         "3",
					WorkflowID: "workflow-2",
					StartTime:  time.Date(2024, 6, 3, 10, 0, 0, 0, time.UTC),
					EndTime:    time.Date(2024, 6, 3, 10, 30, 0, 0, time.UTC),
					Status:     domain.WorkflowStatusCompleted,
				},
			},
			workflowID: "workflow-1",
			expected: []domain.WorkflowRun{
				{
					ID:         "1",
					WorkflowID: "workflow-1",
					StartTime:  time.Date(2024, 6, 1, 12, 0, 0, 0, time.UTC),
					EndTime:    time.Date(2024, 6, 1, 12, 30, 0, 0, time.UTC),
					Status:     domain.WorkflowStatusCompleted,
				},
				{
					ID:         "2",
					WorkflowID: "workflow-1",
					StartTime:  time.Date(2024, 6, 2, 14, 0, 0, 0, time.UTC),
					EndTime:    time.Date(2024, 6, 2, 14, 45, 0, 0, time.UTC),
					Status:     domain.WorkflowStatusFailed,
				},
			},
			wantErr: false,
		},
		{
			name: "failure: no workflow runs for given workflow ID",
			existingRuns: []domain.WorkflowRun{
				{
					ID:         "1",
					WorkflowID: "workflow-1",
					StartTime:  time.Date(2024, 6, 1, 12, 0, 0, 0, time.UTC),
					EndTime:    time.Date(2024, 6, 1, 12, 30, 0, 0, time.UTC),
					Status:     domain.WorkflowStatusCompleted,
				},
			},
			workflowID: "non-existent-workflow-id",
			expected:   []domain.WorkflowRun{},
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
		existingRuns []domain.WorkflowRun
		workflowID   string
		after        time.Time
		expected     []domain.WorkflowRun
		wantErr      bool
	}{
		{
			name: "success: get workflow runs after a specific time",
			existingRuns: []domain.WorkflowRun{
				{
					ID:         "1",
					WorkflowID: "workflow-1",
					StartTime:  time.Date(2024, 6, 1, 12, 0, 0, 0, time.UTC),
					EndTime:    time.Date(2024, 6, 1, 12, 30, 0, 0, time.UTC),
					Status:     domain.WorkflowStatusCompleted,
				},
				{
					ID:         "2",
					WorkflowID: "workflow-1",
					StartTime:  time.Date(2024, 6, 2, 14, 0, 0, 0, time.UTC),
					EndTime:    time.Date(2024, 6, 2, 14, 45, 0, 0, time.UTC),
					Status:     domain.WorkflowStatusFailed,
				},
				{
					ID:         "3",
					WorkflowID: "workflow-1",
					StartTime:  time.Date(2024, 6, 3, 10, 0, 0, 0, time.UTC),
					EndTime:    time.Date(2024, 6, 3, 10, 30, 0, 0, time.UTC),
					Status:     domain.WorkflowStatusCompleted,
				},
			},
			workflowID: "workflow-1",
			after:      time.Date(2024, 6, 2, 0, 0, 0, 0, time.UTC),
			expected: []domain.WorkflowRun{
				{
					ID:         "2",
					WorkflowID: "workflow-1",
					StartTime:  time.Date(2024, 6, 2, 14, 0, 0, 0, time.UTC),
					EndTime:    time.Date(2024, 6, 2, 14, 45, 0, 0, time.UTC),
					Status:     domain.WorkflowStatusFailed,
				},
				{
					ID:         "3",
					WorkflowID: "workflow-1",
					StartTime:  time.Date(2024, 6, 3, 10, 0, 0, 0, time.UTC),
					EndTime:    time.Date(2024, 6, 3, 10, 30, 0, 0, time.UTC),
					Status:     domain.WorkflowStatusCompleted,
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
