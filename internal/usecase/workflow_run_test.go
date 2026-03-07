package usecase

import (
	"errors"
	"os"
	"testing"
	"time"

	"go.uber.org/mock/gomock"
	"github.com/kavos113/seseragi/internal/domain"
	"github.com/kavos113/seseragi/internal/domain/mock_domain"
	"github.com/stretchr/testify/assert"
)

func TestRunWorkflow(t *testing.T) {
	tests := []struct {
		name           string
		workflowID     string
		setupMock      func(wfRepo *mock_domain.MockWorkflowRepository, runRepo *mock_domain.MockWorkflowRunRepository, taskRepo *mock_domain.MockTaskRepository, idGen *testIDGenerator, runner *mock_domain.MockNodeRunner)
		wantErr        bool
	}{
		{
			name:       "success: simple workflow",
			workflowID: "wf-1",
			setupMock: func(wfRepo *mock_domain.MockWorkflowRepository, runRepo *mock_domain.MockWorkflowRunRepository, taskRepo *mock_domain.MockTaskRepository, idGen *testIDGenerator, runner *mock_domain.MockNodeRunner) {
				wf := domain.Workflow{
					ID:   "wf-1",
					Name: "Workflow 1",
					Nodes: []domain.Node{
						{Name: "node-1", TaskName: "task-1"},
					},
				}
				wfRepo.EXPECT().GetWorkflowByID("wf-1").Return(wf, nil)
				idGen.ids = []string{"run-1"}
				
				task := domain.Task{Name: "task-1"}
				taskRepo.EXPECT().GetTaskByName("task-1").Return(task, nil)
				
				runner.EXPECT().Run(wf.Nodes[0], task, "run-1").Return(nil)
				
				runRepo.EXPECT().CreateWorkflowRun(gomock.Any()).DoAndReturn(func(run domain.WorkflowRun) (domain.WorkflowRun, error) {
					assert.Equal(t, "run-1", run.ID)
					assert.Equal(t, domain.WorkflowStatusCompleted, run.Status)
					return run, nil
				})
			},
			wantErr: false,
		},
		{
			name:       "success: workflow with dependencies",
			workflowID: "wf-1",
			setupMock: func(wfRepo *mock_domain.MockWorkflowRepository, runRepo *mock_domain.MockWorkflowRunRepository, taskRepo *mock_domain.MockTaskRepository, idGen *testIDGenerator, runner *mock_domain.MockNodeRunner) {
				wf := domain.Workflow{
					ID:   "wf-1",
					Nodes: []domain.Node{
						{Name: "node-1", TaskName: "task-1"},
						{Name: "node-2", TaskName: "task-2", Dependencies: []string{"node-1"}},
					},
				}
				wfRepo.EXPECT().GetWorkflowByID("wf-1").Return(wf, nil)
				idGen.ids = []string{"run-1"}
				
				task1 := domain.Task{Name: "task-1"}
				task2 := domain.Task{Name: "task-2"}
				taskRepo.EXPECT().GetTaskByName("task-1").Return(task1, nil)
				taskRepo.EXPECT().GetTaskByName("task-2").Return(task2, nil)
				
				runner.EXPECT().Run(wf.Nodes[0], task1, "run-1").Return(nil)
				runner.EXPECT().Run(wf.Nodes[1], task2, "run-1").Return(nil)
				
				runRepo.EXPECT().CreateWorkflowRun(gomock.Any()).Return(domain.WorkflowRun{}, nil)
			},
			wantErr: false,
		},
		{
			name:       "failure: node error",
			workflowID: "wf-1",
			setupMock: func(wfRepo *mock_domain.MockWorkflowRepository, runRepo *mock_domain.MockWorkflowRunRepository, taskRepo *mock_domain.MockTaskRepository, idGen *testIDGenerator, runner *mock_domain.MockNodeRunner) {
				wf := domain.Workflow{
					ID:   "wf-1",
					Nodes: []domain.Node{
						{Name: "node-1", TaskName: "task-1"},
					},
				}
				wfRepo.EXPECT().GetWorkflowByID("wf-1").Return(wf, nil)
				idGen.ids = []string{"run-1"}
				
				task := domain.Task{Name: "task-1"}
				taskRepo.EXPECT().GetTaskByName("task-1").Return(task, nil)
				
				runner.EXPECT().Run(wf.Nodes[0], task, "run-1").Return(errors.New("node error"))
				
				runRepo.EXPECT().CreateWorkflowRun(gomock.Any()).DoAndReturn(func(run domain.WorkflowRun) (domain.WorkflowRun, error) {
					assert.Equal(t, domain.WorkflowStatusFailed, run.Status)
					return run, nil
				})
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			wfRepo := mock_domain.NewMockWorkflowRepository(ctrl)
			runRepo := mock_domain.NewMockWorkflowRunRepository(ctrl)
			taskRepo := mock_domain.NewMockTaskRepository(ctrl)
			idGen := newTestIDGenerator()
			runner := mock_domain.NewMockNodeRunner(ctrl)

			tt.setupMock(wfRepo, runRepo, taskRepo, idGen, runner)

			uc := NewWorkflowRunUseCase(wfRepo, runRepo, taskRepo, idGen)
			err := uc.RunWorkflow(tt.workflowID, func(node domain.Node) domain.NodeRunner {
				return runner
			})

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			
			// Cleanup data dir created by uc.RunWorkflow
			os.RemoveAll(domain.GetDataDir("run-1"))
		})
	}
}

func TestGetWorkflowsToRun(t *testing.T) {
	tests := []struct {
		name      string
		setupMock func(wfRepo *mock_domain.MockWorkflowRepository, runRepo *mock_domain.MockWorkflowRunRepository)
		expected  int
	}{
		{
			name: "run all",
			setupMock: func(wfRepo *mock_domain.MockWorkflowRepository, runRepo *mock_domain.MockWorkflowRunRepository) {
				wfs := []domain.Workflow{
					{ID: "wf-1", RunInterval: time.Hour},
					{ID: "wf-2", RunInterval: time.Hour},
				}
				wfRepo.EXPECT().GetAllWorkflows().Return(wfs, nil)
				runRepo.EXPECT().GetWorkflowRunsAfter("wf-1", gomock.Any()).Return([]domain.WorkflowRun{}, nil)
				runRepo.EXPECT().GetWorkflowRunsAfter("wf-2", gomock.Any()).Return([]domain.WorkflowRun{}, nil)
			},
			expected: 2,
		},
		{
			name: "run one, skip one",
			setupMock: func(wfRepo *mock_domain.MockWorkflowRepository, runRepo *mock_domain.MockWorkflowRunRepository) {
				wfs := []domain.Workflow{
					{ID: "wf-1", RunInterval: time.Hour},
					{ID: "wf-2", RunInterval: time.Hour},
				}
				wfRepo.EXPECT().GetAllWorkflows().Return(wfs, nil)
				runRepo.EXPECT().GetWorkflowRunsAfter("wf-1", gomock.Any()).Return([]domain.WorkflowRun{{ID: "run-1"}}, nil)
				runRepo.EXPECT().GetWorkflowRunsAfter("wf-2", gomock.Any()).Return([]domain.WorkflowRun{}, nil)
			},
			expected: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			wfRepo := mock_domain.NewMockWorkflowRepository(ctrl)
			runRepo := mock_domain.NewMockWorkflowRunRepository(ctrl)

			tt.setupMock(wfRepo, runRepo)

			uc := NewWorkflowRunUseCase(wfRepo, runRepo, nil, nil)
			wfs, err := uc.GetWorkflowsToRun()

			assert.NoError(t, err)
			assert.Len(t, wfs, tt.expected)
		})
	}
}
