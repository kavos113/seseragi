package manager

import (
	"slices"
	"testing"

	"github.com/kavos113/seseragi/model"
	"github.com/kavos113/seseragi/model/mock_model"
	"go.uber.org/mock/gomock"
)

func TestParseWorkflow(t *testing.T) {
	tests := []struct {
		name         string
		workflowInfo *WorkflowInfo
		setupMock    func(repo *mock_model.MockTaskRepository)
		want         model.Workflow
		wantErr      bool
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
					{TaskID: "task-id-1", Dependencies: []*model.Node{}},
				},
			},
			wantErr: false,
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
					{TaskID: "task-id-1", Dependencies: []*model.Node{{TaskID: "task-id-2"}}},
					{TaskID: "task-id-2", Dependencies: []*model.Node{}},
				},
			},
			wantErr: false,
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
					{TaskID: "task-id-1", Dependencies: []*model.Node{{TaskID: "task-id-2"}, {TaskID: "task-id-3"}}},
					{TaskID: "task-id-2", Dependencies: []*model.Node{}},
					{TaskID: "task-id-3", Dependencies: []*model.Node{}},
				},
			},
			wantErr: false,
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
					{TaskID: "task-id-1", Dependencies: []*model.Node{{TaskID: "task-id-2"}, {TaskID: "task-id-3"}}},
					{TaskID: "task-id-3", Dependencies: []*model.Node{{TaskID: "task-id-2"}}},
					{TaskID: "task-id-2", Dependencies: []*model.Node{}},
				},
			},
			wantErr: false,
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
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockTaskRepo := mock_model.NewMockTaskRepository(ctrl)
			tt.setupMock(mockTaskRepo)
			taskRepo = mockTaskRepo

			got, err := ParseWorkflow(tt.workflowInfo, "")
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseWorkflow() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if got.Name != tt.want.Name {
				t.Errorf("ParseWorkflow() got.Name = %v, want %v", got.Name, tt.want.Name)
			}
			if len(got.Nodes) != len(tt.want.Nodes) {
				t.Errorf("ParseWorkflow() got.Nodes length = %v, want %v", len(got.Nodes), len(tt.want.Nodes))
				return
			}

			for _, gotNode := range got.Nodes {
				wantNodeIndex := slices.IndexFunc(tt.want.Nodes, func(n model.Node) bool {
					return n.TaskID == gotNode.TaskID
				})
				if wantNodeIndex == -1 {
					t.Errorf("ParseWorkflow() got.Nodes[%v] has unexpected TaskID = %v", wantNodeIndex, gotNode.TaskID)
					continue
				}
				wantNode := tt.want.Nodes[wantNodeIndex]

				if len(gotNode.Dependencies) != len(wantNode.Dependencies) {
					t.Errorf("ParseWorkflow() got.Nodes[%v].Dependencies length = %v, want %v", wantNodeIndex, len(gotNode.Dependencies), len(wantNode.Dependencies))
					continue
				}

				for j, gotDep := range gotNode.Dependencies {
					wantDepIndex := slices.IndexFunc(wantNode.Dependencies, func(d *model.Node) bool {
						return d.TaskID == gotDep.TaskID
					})
					if wantDepIndex == -1 {
						t.Errorf("ParseWorkflow() got.Nodes[%v].Dependencies[%v] has unexpected TaskID = %v", wantNodeIndex, j, gotDep.TaskID)
					}
				}
			}
		})
	}
}
