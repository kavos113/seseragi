package usecase

import (
	"testing"

	"github.com/kavos113/seseragi/internal/domain"
	"github.com/kavos113/seseragi/internal/domain/mock_domain"
	"go.uber.org/mock/gomock"
)

func TestAddTask(t *testing.T) {
	tests := []struct {
		name      string
		task      domain.Task
		setupMock func(repo *mock_domain.MockTaskRepository, provider *mock_domain.MockTaskProvider)
		wantErr   bool
	}{
		{
			name: "success",
			task: domain.Task{Name: "Test Task"},
			setupMock: func(repo *mock_domain.MockTaskRepository, provider *mock_domain.MockTaskProvider) {
				repo.EXPECT().
					CreateTask(domain.Task{
						ID:   "generated-id",
						Name: "Test Task",
					}).
					Return(domain.Task{
						ID:   "generated-id",
						Name: "Test Task",
					}, nil)
				provider.EXPECT().
					BuildTask(domain.Task{
						ID:   "generated-id",
						Name: "Test Task",
					}).
					Return(nil)
			},
			wantErr: false,
		},
		{
			name: "success: docker task fills image name",
			task: domain.Task{
				Name: "Docker Task",
				TaskDef: domain.DockerTaskDefinition{
					ContextDir: "./context",
				},
			},
			setupMock: func(repo *mock_domain.MockTaskRepository, provider *mock_domain.MockTaskProvider) {
				expectedTask := domain.Task{
					ID:   "generated-id",
					Name: "Docker Task",
					TaskDef: domain.DockerTaskDefinition{
						ImageName:  "generated-id",
						ContextDir: "./context",
					},
				}
				repo.EXPECT().
					CreateTask(expectedTask).
					Return(expectedTask, nil)
				provider.EXPECT().
					BuildTask(expectedTask).
					Return(nil)
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockRepo := mock_domain.NewMockTaskRepository(ctrl)
			mockProvider := mock_domain.NewMockTaskProvider(ctrl)
			mockIDGen := newTestIDGenerator("generated-id")

			tt.setupMock(mockRepo, mockProvider)

			uc := NewTaskUseCase(mockRepo, mockProvider, mockIDGen)
			err := uc.AddTask(tt.task)
			if err != nil {
				if !tt.wantErr {
					t.Errorf("AddTask() error = %v, wantErr %v", err, tt.wantErr)
				}
				return
			}
		})
	}
}
