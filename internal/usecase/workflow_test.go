package usecase

import (
	"testing"

	"github.com/kavos113/seseragi/internal/domain"
	"github.com/kavos113/seseragi/internal/domain/mock_domain"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestAddWorkflow(t *testing.T) {
	tests := []struct {
		name         string
		workflow     domain.Workflow
		setupMock    func(repo *mock_domain.MockWorkflowRepository, taskRepo *mock_domain.MockTaskRepository)
		wantWorkflow domain.Workflow
		wantErr      error
	}{
		{
			name: "success: simple workflow",
			workflow: domain.Workflow{
				Name: "hello workflow",
				Nodes: []domain.Node{
					{Name: "node1", TaskName: "task1", Dependencies: []string{}},
				},
			},
			setupMock: func(repo *mock_domain.MockWorkflowRepository, taskRepo *mock_domain.MockTaskRepository) {
				repo.EXPECT().
					CreateWorkflow(domain.Workflow{
						ID:   "1",
						Name: "hello workflow",
						Nodes: []domain.Node{
							{Name: "node1", TaskName: "task1", Dependencies: []string{}},
						},
					}).
					Return(domain.Workflow{
						ID:   "1",
						Name: "hello workflow",
						Nodes: []domain.Node{
							{Name: "node1", TaskName: "task1", Dependencies: []string{}},
						},
					}, nil)
				taskRepo.EXPECT().
					GetTaskByName("task1").
					Return(domain.Task{Name: "task1"}, nil)
			},
			wantWorkflow: domain.Workflow{
				ID:   "1",
				Name: "hello workflow",
				Nodes: []domain.Node{
					{Name: "node1", TaskName: "task1", Dependencies: []string{}},
				},
			},
			wantErr: nil,
		},
		{
			name: "success: workflow with dependencies",
			workflow: domain.Workflow{
				Name: "hello workflow",
				Nodes: []domain.Node{
					{Name: "node1", TaskName: "task1", Dependencies: []string{}},
					{Name: "node2", TaskName: "task2", Dependencies: []string{"node1"}},
				},
			},
			setupMock: func(repo *mock_domain.MockWorkflowRepository, taskRepo *mock_domain.MockTaskRepository) {
				repo.EXPECT().
					CreateWorkflow(domain.Workflow{
						ID:   "1",
						Name: "hello workflow",
						Nodes: []domain.Node{
							{Name: "node1", TaskName: "task1", Dependencies: []string{}},
							{Name: "node2", TaskName: "task2", Dependencies: []string{"node1"}},
						},
					}).
					Return(domain.Workflow{
						ID:   "1",
						Name: "hello workflow",
						Nodes: []domain.Node{
							{Name: "node1", TaskName: "task1", Dependencies: []string{}},
							{Name: "node2", TaskName: "task2", Dependencies: []string{"node1"}},
						},
					}, nil)
				taskRepo.EXPECT().
					GetTaskByName("task1").
					Return(domain.Task{Name: "task1"}, nil)
				taskRepo.EXPECT().
					GetTaskByName("task2").
					Return(domain.Task{Name: "task2"}, nil)
			},
			wantWorkflow: domain.Workflow{
				ID:   "1",
				Name: "hello workflow",
				Nodes: []domain.Node{
					{Name: "node1", TaskName: "task1", Dependencies: []string{}},
					{Name: "node2", TaskName: "task2", Dependencies: []string{"node1"}},
				},
			},
			wantErr: nil,
		},
		{
			name: "success: multiple dependencies",
			workflow: domain.Workflow{
				Name: "hello workflow",
				Nodes: []domain.Node{
					{Name: "node1", TaskName: "task1", Dependencies: []string{}},
					{Name: "node2", TaskName: "task2", Dependencies: []string{"node1"}},
					{Name: "node3", TaskName: "task3", Dependencies: []string{"node1"}},
					{Name: "node4", TaskName: "task4", Dependencies: []string{"node2", "node3"}},
				},
			},
			setupMock: func(repo *mock_domain.MockWorkflowRepository, taskRepo *mock_domain.MockTaskRepository) {
				repo.EXPECT().
					CreateWorkflow(domain.Workflow{
						ID:   "1",
						Name: "hello workflow",
						Nodes: []domain.Node{
							{Name: "node1", TaskName: "task1", Dependencies: []string{}},
							{Name: "node2", TaskName: "task2", Dependencies: []string{"node1"}},
							{Name: "node3", TaskName: "task3", Dependencies: []string{"node1"}},
							{Name: "node4", TaskName: "task4", Dependencies: []string{"node2", "node3"}},
						},
					}).
					Return(domain.Workflow{
						ID:   "1",
						Name: "hello workflow",
						Nodes: []domain.Node{
							{Name: "node1", TaskName: "task1", Dependencies: []string{}},
							{Name: "node2", TaskName: "task2", Dependencies: []string{"node1"}},
							{Name: "node3", TaskName: "task3", Dependencies: []string{"node1"}},
							{Name: "node4", TaskName: "task4", Dependencies: []string{"node2", "node3"}},
						},
					}, nil)
				taskRepo.EXPECT().
					GetTaskByName("task1").
					Return(domain.Task{Name: "task1"}, nil)
				taskRepo.EXPECT().
					GetTaskByName("task2").
					Return(domain.Task{Name: "task2"}, nil)
				taskRepo.EXPECT().
					GetTaskByName("task3").
					Return(domain.Task{Name: "task3"}, nil)
				taskRepo.EXPECT().
					GetTaskByName("task4").
					Return(domain.Task{Name: "task4"}, nil)
			},
			wantWorkflow: domain.Workflow{
				ID:   "1",
				Name: "hello workflow",
				Nodes: []domain.Node{
					{Name: "node1", TaskName: "task1", Dependencies: []string{}},
					{Name: "node2", TaskName: "task2", Dependencies: []string{"node1"}},
					{Name: "node3", TaskName: "task3", Dependencies: []string{"node1"}},
					{Name: "node4", TaskName: "task4", Dependencies: []string{"node2", "node3"}},
				},
			},
			wantErr: nil,
		},
		{
			name: "failure: circular dependency",
			workflow: domain.Workflow{
				Name: "circular workflow",
				Nodes: []domain.Node{
					{Name: "node1", TaskName: "task1", Dependencies: []string{"node3"}},
					{Name: "node2", TaskName: "task2", Dependencies: []string{"node1"}},
					{Name: "node3", TaskName: "task3", Dependencies: []string{"node2"}},
				},
			},
			setupMock:    func(repo *mock_domain.MockWorkflowRepository, taskRepo *mock_domain.MockTaskRepository) {},
			wantWorkflow: domain.Workflow{},
			wantErr:      ErrWorkflowCircularDependency,
		},
		{
			name: "failure: missing dependency",
			workflow: domain.Workflow{
				Name: "missing dependency workflow",
				Nodes: []domain.Node{
					{Name: "node1", TaskName: "task1", Dependencies: []string{"node2"}},
				},
			},
			setupMock:    func(repo *mock_domain.MockWorkflowRepository, taskRepo *mock_domain.MockTaskRepository) {},
			wantWorkflow: domain.Workflow{},
			wantErr:      ErrWorkflowMissingDependency,
		},
		{
			name: "failure: missing task",
			workflow: domain.Workflow{
				Name: "missing task workflow",
				Nodes: []domain.Node{
					{Name: "node1", TaskName: "task1", Dependencies: []string{}},
				},
			},
			setupMock: func(repo *mock_domain.MockWorkflowRepository, taskRepo *mock_domain.MockTaskRepository) {
				taskRepo.EXPECT().
					GetTaskByName("task1").
					Return(domain.Task{}, assert.AnError)
			},
			wantWorkflow: domain.Workflow{},
			wantErr:      ErrWorkflowMissingTask,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockRepo := mock_domain.NewMockWorkflowRepository(ctrl)
			mockTaskRepo := mock_domain.NewMockTaskRepository(ctrl)
			mockIDProvider := newTestIDGenerator("1", "2", "3", "4", "5")
			tt.setupMock(mockRepo, mockTaskRepo)

			uc := NewWorkflowUseCase(mockRepo, mockTaskRepo, mockIDProvider)
			got, err := uc.AddWorkflow(tt.workflow)
			if err != nil {
				assert.ErrorIs(t, err, tt.wantErr)
				return
			}

			assert.Equal(t, tt.wantWorkflow, got)
		})
	}
}
