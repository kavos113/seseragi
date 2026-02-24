package manager

import (
	"errors"
	"testing"

	"github.com/kavos113/seseragi/model"
	"github.com/kavos113/seseragi/model/mock_model"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestParseWorkflow(t *testing.T) {
	tests := []struct {
		name         string
		workflowInfo *WorkflowInfo
		setupMock    func(repo *mock_model.MockTaskRepository)
		want         model.Workflow
		wantErr      error
	}{
		{
			name: "success: simple workflow",
			workflowInfo: &WorkflowInfo{
				Name:        "hello-workflow",
				Description: "Hello Workflow",
				Nodes: map[string]NodeInfo{
					"go-hello": {ID: "task-id-1"},
				},
			},
			setupMock: func(repo *mock_model.MockTaskRepository) {
				repo.EXPECT().
					GetTaskByID("task-id-1").
					Return(model.Task{ID: "task-id-1"}, nil)
			},
			want: model.Workflow{
				Name: "hello-workflow",
				Nodes: []model.Node{
					{TaskID: "task-id-1", Dependencies: []string{}},
				},
			},
			wantErr: nil,
		},
		{
			name: "success: workflow with dependencies",
			workflowInfo: &WorkflowInfo{
				Name:        "hello-workflow",
				Description: "Hello Workflow",
				Nodes: map[string]NodeInfo{
					"go-hello": {ID: "task-id-1", Dependencies: []string{"go-world"}},
					"go-world": {ID: "task-id-2"},
				},
			},
			setupMock: func(repo *mock_model.MockTaskRepository) {
				repo.EXPECT().
					GetTaskByID("task-id-1").
					Return(model.Task{ID: "task-id-1"}, nil)
				repo.EXPECT().
					GetTaskByID("task-id-2").
					Return(model.Task{ID: "task-id-2"}, nil)
			},
			want: model.Workflow{
				Name: "hello-workflow",
				Nodes: []model.Node{
					{TaskID: "task-id-1", Dependencies: []string{"task-id-2"}},
					{TaskID: "task-id-2", Dependencies: []string{}},
				},
			},
			wantErr: nil,
		},
		{
			name: "success: multiple dependencies",
			workflowInfo: &WorkflowInfo{
				Name:        "hello-workflow",
				Description: "Hello Workflow",
				Nodes: map[string]NodeInfo{
					"go-hello":    {ID: "task-id-1", Dependencies: []string{"go-world", "go-universe"}},
					"go-world":    {ID: "task-id-2"},
					"go-universe": {ID: "task-id-3"},
				},
			},
			setupMock: func(repo *mock_model.MockTaskRepository) {
				repo.EXPECT().
					GetTaskByID("task-id-1").
					Return(model.Task{ID: "task-id-1"}, nil)
				repo.EXPECT().
					GetTaskByID("task-id-2").
					Return(model.Task{ID: "task-id-2"}, nil)
				repo.EXPECT().
					GetTaskByID("task-id-3").
					Return(model.Task{ID: "task-id-3"}, nil)
			},
			want: model.Workflow{
				Name: "hello-workflow",
				Nodes: []model.Node{
					{TaskID: "task-id-1", Dependencies: []string{"task-id-2", "task-id-3"}},
					{TaskID: "task-id-2", Dependencies: []string{}},
					{TaskID: "task-id-3", Dependencies: []string{}},
				},
			},
			wantErr: nil,
		},
		{
			name: "success: shared dependencies",
			workflowInfo: &WorkflowInfo{
				Name:        "hello-workflow",
				Description: "Hello Workflow",
				Nodes: map[string]NodeInfo{
					"go-hello":    {ID: "task-id-1", Dependencies: []string{"go-world", "go-universe"}},
					"go-universe": {ID: "task-id-3", Dependencies: []string{"go-world"}},
					"go-world":    {ID: "task-id-2"},
				},
			},
			setupMock: func(repo *mock_model.MockTaskRepository) {
				repo.EXPECT().
					GetTaskByID("task-id-1").
					Return(model.Task{ID: "task-id-1"}, nil)
				repo.EXPECT().
					GetTaskByID("task-id-2").
					Return(model.Task{ID: "task-id-2"}, nil)
				repo.EXPECT().
					GetTaskByID("task-id-3").
					Return(model.Task{ID: "task-id-3"}, nil)
			},
			want: model.Workflow{
				Name: "hello-workflow",
				Nodes: []model.Node{
					{TaskID: "task-id-1", Dependencies: []string{"task-id-2", "task-id-3"}},
					{TaskID: "task-id-3", Dependencies: []string{"task-id-2"}},
					{TaskID: "task-id-2", Dependencies: []string{}},
				},
			},
			wantErr: nil,
		},
		{
			name: "failure: missing task for node",
			workflowInfo: &WorkflowInfo{
				Name:        "hello-workflow",
				Description: "Hello Workflow",
				Nodes: map[string]NodeInfo{
					"go-hello": {ID: "task-id-1", Dependencies: []string{"go-world"}},
				},
			},
			setupMock: func(repo *mock_model.MockTaskRepository) {
				repo.EXPECT().
					GetTaskByID("task-id-1").
					Return(model.Task{ID: "task-id-1"}, nil)
			},
			want:    model.Workflow{},
			wantErr: assert.AnError,
		},
		{
			name: "failure: circular dependency",
			workflowInfo: &WorkflowInfo{
				Name:        "hello-workflow",
				Description: "Hello Workflow",
				Nodes: map[string]NodeInfo{
					"go-hello": {ID: "task-id-1", Dependencies: []string{"go-world"}},
					"go-world": {ID: "task-id-2", Dependencies: []string{"go-hello"}},
				},
			},
			setupMock: func(repo *mock_model.MockTaskRepository) {
				repo.EXPECT().
					GetTaskByID("task-id-1").
					Return(model.Task{ID: "task-id-1"}, nil)
				repo.EXPECT().
					GetTaskByID("task-id-2").
					Return(model.Task{ID: "task-id-2"}, nil)
			},
			want:    model.Workflow{},
			wantErr: ErrWorkflowCircularDependency,
		},
		{
			name: "failure: circular dependency with multiple nodes",
			workflowInfo: &WorkflowInfo{
				Name:        "hello-workflow",
				Description: "Hello Workflow",
				Nodes: map[string]NodeInfo{
					"go-hello":    {ID: "task-id-1", Dependencies: []string{"go-world"}},
					"go-world":    {ID: "task-id-2", Dependencies: []string{"go-universe"}},
					"go-universe": {ID: "task-id-3", Dependencies: []string{"go-hello"}},
				},
			},
			setupMock: func(repo *mock_model.MockTaskRepository) {
				repo.EXPECT().
					GetTaskByID("task-id-1").
					Return(model.Task{ID: "task-id-1"}, nil)
				repo.EXPECT().
					GetTaskByID("task-id-2").
					Return(model.Task{ID: "task-id-2"}, nil)
				repo.EXPECT().
					GetTaskByID("task-id-3").
					Return(model.Task{ID: "task-id-3"}, nil)
			},
			want:    model.Workflow{},
			wantErr: ErrWorkflowCircularDependency,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			oldRepo := taskRepo
			t.Cleanup(func() {
				taskRepo = oldRepo
			})

			mockTaskRepo := mock_model.NewMockTaskRepository(ctrl)
			tt.setupMock(mockTaskRepo)
			taskRepo = mockTaskRepo

			got, err := ParseWorkflow(tt.workflowInfo, "")
			if tt.wantErr != nil {
				if errors.Is(tt.wantErr, assert.AnError) {
					assert.Error(t, err)
				} else {
					assert.ErrorIs(t, err, tt.wantErr)
				}
				return
			}
			assert.NoError(t, err)

			assert.Equal(t, got.Name, tt.want.Name)
			assert.Equal(t, got.YamlPath, tt.want.YamlPath)
			assert.Equal(t, len(got.Nodes), len(tt.want.Nodes))
			assert.ElementsMatch(t, got.Nodes, tt.want.Nodes)
		})
	}
}
