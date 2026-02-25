package usecase

import (
	"testing"

	"github.com/kavos113/seseragi/internal/domain"
	"github.com/kavos113/seseragi/internal/domain/mock_domain"
	"github.com/stretchr/testify/assert"
)

func TestAddWorkflow(t *testing.T) {
	tests := []struct {
		name         string
		workflow     domain.Workflow
		setupMock    func(repo *mock_domain.MockWorkflowRepository)
		wantWorkflow domain.Workflow
		wantErr      error
	}{
		{
			name: "success: simple workflow",
			workflow: domain.Workflow{
				Name: "hello workflow",
				Nodes: []domain.Node{
					{Name: "node1", Dependencies: []string{}},
				},
			},
			setupMock: func(repo *mock_domain.MockWorkflowRepository) {
				repo.EXPECT().
					CreateWorkflow(domain.Workflow{
						ID:   "1",
						Name: "hello workflow",
						Nodes: []domain.Node{
							{Name: "node1", Dependencies: []string{}},
						},
					}).
					Return(domain.Workflow{
						ID:   "1",
						Name: "hello workflow",
						Nodes: []domain.Node{
							{Name: "node1", Dependencies: []string{}},
						},
					}, nil)
			},
			wantWorkflow: domain.Workflow{
				ID:   "1",
				Name: "hello workflow",
				Nodes: []domain.Node{
					{Name: "node1", Dependencies: []string{}},
				},
			},
			wantErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := &mock_domain.MockWorkflowRepository{}
			mockIDProvider := newTestIDGenerator("1", "2", "3", "4", "5")
			tt.setupMock(mockRepo)

			uc := NewWorkflowUseCase(mockRepo, mockIDProvider)
			err := uc.AddWorkflow(tt.workflow)
			if err != nil {
				assert.ErrorIs(t, err, tt.wantErr)
				return
			}
		})
	}
}
