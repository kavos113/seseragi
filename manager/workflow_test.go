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
					"go-hello": {Name: "go-hello", Dependencies: []string{}},
				},
			},
			setupMock: func(repo *mock_model.MockTaskRepository) {
				repo.EXPECT().
					GetTaskByName("go-hello").
					Return(model.Task{ID: "task-id-1"}, nil)
			},
			want: model.Workflow{
				Name: "hello-workflow",
				Nodes: []model.Node{
					{Name: "go-hello", TaskID: "task-id-1", Dependencies: []string{}},
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
					"go-hello": {Name: "go-hello", Dependencies: []string{"go-world"}},
					"go-world": {Name: "go-world", Dependencies: []string{}},
				},
			},
			setupMock: func(repo *mock_model.MockTaskRepository) {
				repo.EXPECT().
					GetTaskByName("go-hello").
					Return(model.Task{ID: "task-id-1"}, nil)
				repo.EXPECT().
					GetTaskByName("go-world").
					Return(model.Task{ID: "task-id-2"}, nil)
			},
			want: model.Workflow{
				Name: "hello-workflow",
				Nodes: []model.Node{
					{Name: "go-hello", TaskID: "task-id-1", Dependencies: []string{"go-world"}},
					{Name: "go-world", TaskID: "task-id-2", Dependencies: []string{}},
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
					"go-hello":    {Name: "go-hello", Dependencies: []string{"go-world", "go-universe"}},
					"go-world":    {Name: "go-world"},
					"go-universe": {Name: "go-universe"},
				},
			},
			setupMock: func(repo *mock_model.MockTaskRepository) {
				repo.EXPECT().
					GetTaskByName("go-hello").
					Return(model.Task{ID: "task-id-1"}, nil)
				repo.EXPECT().
					GetTaskByName("go-world").
					Return(model.Task{ID: "task-id-2"}, nil)
				repo.EXPECT().
					GetTaskByName("go-universe").
					Return(model.Task{ID: "task-id-3"}, nil)
			},
			want: model.Workflow{
				Name: "hello-workflow",
				Nodes: []model.Node{
					{Name: "go-hello", TaskID: "task-id-1", Dependencies: []string{"go-world", "go-universe"}},
					{Name: "go-world", TaskID: "task-id-2", Dependencies: []string{}},
					{Name: "go-universe", TaskID: "task-id-3", Dependencies: []string{}},
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
					"go-hello":    {Name: "go-hello", Dependencies: []string{"go-world", "go-universe"}},
					"go-universe": {Name: "go-universe", Dependencies: []string{"go-world"}},
					"go-world":    {Name: "go-world"},
				},
			},
			setupMock: func(repo *mock_model.MockTaskRepository) {
				repo.EXPECT().
					GetTaskByName("go-hello").
					Return(model.Task{ID: "task-id-1"}, nil)
				repo.EXPECT().
					GetTaskByName("go-universe").
					Return(model.Task{ID: "task-id-2"}, nil)
				repo.EXPECT().
					GetTaskByName("go-world").
					Return(model.Task{ID: "task-id-3"}, nil)
			},
			want: model.Workflow{
				Name: "hello-workflow",
				Nodes: []model.Node{
					{Name: "go-hello", TaskID: "task-id-1", Dependencies: []string{"go-world", "go-universe"}},
					{Name: "go-universe", TaskID: "task-id-2", Dependencies: []string{"go-world"}},
					{Name: "go-world", TaskID: "task-id-3", Dependencies: []string{}},
				},
			},
			wantErr: nil,
		},
		{
			name: "failure: missing dependency",
			workflowInfo: &WorkflowInfo{
				Name:        "hello-workflow",
				Description: "Hello Workflow",
				Nodes: map[string]NodeInfo{
					"go-hello": {Name: "go-hello", Dependencies: []string{"go-world"}},
				},
			},
			setupMock: func(repo *mock_model.MockTaskRepository) {
				repo.EXPECT().
					GetTaskByName("go-hello").
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
					"go-hello": {Name: "go-hello", Dependencies: []string{"go-world"}},
					"go-world": {Name: "go-world", Dependencies: []string{"go-hello"}},
				},
			},
			setupMock: func(repo *mock_model.MockTaskRepository) {
				repo.EXPECT().
					GetTaskByName("go-hello").
					Return(model.Task{ID: "task-id-1"}, nil)
				repo.EXPECT().
					GetTaskByName("go-world").
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
					"go-hello":    {Name: "go-hello", Dependencies: []string{"go-world"}},
					"go-world":    {Name: "go-world", Dependencies: []string{"go-universe"}},
					"go-universe": {Name: "go-universe", Dependencies: []string{"go-hello"}},
				},
			},
			setupMock: func(repo *mock_model.MockTaskRepository) {
				repo.EXPECT().
					GetTaskByName("go-hello").
					Return(model.Task{ID: "task-id-1"}, nil)
				repo.EXPECT().
					GetTaskByName("go-world").
					Return(model.Task{ID: "task-id-2"}, nil)
				repo.EXPECT().
					GetTaskByName("go-universe").
					Return(model.Task{ID: "task-id-3"}, nil)
			},
			want:    model.Workflow{},
			wantErr: ErrWorkflowCircularDependency,
		},
		{
			name: "failure: task not found",
			workflowInfo: &WorkflowInfo{
				Name:        "hello-workflow",
				Description: "Hello Workflow",
				Nodes: map[string]NodeInfo{
					"go-hello": {Name: "go-hello", Dependencies: []string{}},
				},
			},
			setupMock: func(repo *mock_model.MockTaskRepository) {
				repo.EXPECT().
					GetTaskByName("go-hello").
					Return(model.Task{}, errors.New("task not found"))
			},
			want:    model.Workflow{},
			wantErr: assert.AnError,
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
