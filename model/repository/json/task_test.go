package json

import (
	"testing"

	"github.com/kavos113/seseragi/model"
	"github.com/stretchr/testify/assert"
)

func TestCreateTask(t *testing.T) {
	tests := []struct {
		name     string
		task     model.Task
		existing []model.Task
		expected model.Task
		wantErr  bool
	}{
		{
			name: "success: create new task",
			task: model.Task{
				ID:   "1",
				Name: "Task 1",
			},
			existing: nil,
			expected: model.Task{
				ID:   "1",
				Name: "Task 1",
			},
			wantErr: false,
		},
		{
			name: "failure: task with same ID already exists",
			task: model.Task{
				ID:   "1",
				Name: "Task 1",
			},
			existing: []model.Task{
				{
					ID:   "1",
					Name: "Existing Task",
				},
			},
			expected: model.Task{},
			wantErr:  true,
		},
		{
			name: "failure: task with same name already exists",
			task: model.Task{
				ID:   "2",
				Name: "Existing Task",
			},
			existing: []model.Task{
				{
					ID:   "1",
					Name: "Existing Task",
				},
			},
			expected: model.Task{},
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := setupTestRepository(t)
			taskRepo := NewJSONTaskRepository(repo)

			for _, existingTask := range tt.existing {
				_, err := taskRepo.CreateTask(existingTask)
				if err != nil {
					t.Fatalf("Failed to create existing task: %v", err)
				}
			}

			created, err := taskRepo.CreateTask(tt.task)
			if err != nil {
				if !tt.wantErr {
					t.Errorf("CreateTask() error = %v, wantErr %v", err, tt.wantErr)
				}
				return
			}

			assert.Equal(t, tt.expected, created)
		})
	}
}

func TestGetTaskByID(t *testing.T) {
	tests := []struct {
		name     string
		existing []model.Task
		id       string
		expected model.Task
		wantErr  bool
	}{
		{
			name: "success: get existing task by ID",
			existing: []model.Task{
				{
					ID:   "1",
					Name: "Task 1",
				},
			},
			id: "1",
			expected: model.Task{
				ID:   "1",
				Name: "Task 1",
			},
			wantErr: false,
		},
		{
			name: "failure: task not found by ID",
			existing: []model.Task{
				{
					ID:   "1",
					Name: "Task 1",
				},
			},
			id:       "non-existent-id",
			expected: model.Task{},
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := setupTestRepository(t)
			taskRepo := NewJSONTaskRepository(repo)

			for _, existingTask := range tt.existing {
				_, err := taskRepo.CreateTask(existingTask)
				if err != nil {
					t.Fatalf("Failed to create existing task: %v", err)
				}
			}

			got, err := taskRepo.GetTaskByID(tt.id)
			if err != nil {
				if !tt.wantErr {
					t.Errorf("GetTaskByID() error = %v, wantErr %v", err, tt.wantErr)
				}
				return
			}

			assert.Equal(t, tt.expected, got)
		})
	}
}

func TestGetTaskByName(t *testing.T) {
	tests := []struct {
		name      string
		existing  []model.Task
		nameToGet string
		expected  model.Task
		wantErr   bool
	}{
		{
			name: "success: get existing task by name",
			existing: []model.Task{
				{
					ID:   "1",
					Name: "Task 1",
				},
			},
			nameToGet: "Task 1",
			expected: model.Task{
				ID:   "1",
				Name: "Task 1",
			},
			wantErr: false,
		},
		{
			name: "failure: task not found by name",
			existing: []model.Task{
				{
					ID:   "1",
					Name: "Task 1",
				},
			},
			nameToGet: "Non-existent Task",
			expected:  model.Task{},
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := setupTestRepository(t)
			taskRepo := NewJSONTaskRepository(repo)

			for _, existingTask := range tt.existing {
				_, err := taskRepo.CreateTask(existingTask)
				if err != nil {
					t.Fatalf("Failed to create existing task: %v", err)
				}
			}

			got, err := taskRepo.GetTaskByName(tt.nameToGet)
			if err != nil {
				if !tt.wantErr {
					t.Errorf("GetTaskByName() error = %v, wantErr %v", err, tt.wantErr)
				}
				return
			}
			assert.Equal(t, tt.expected, got)
		})
	}
}

func TestGetAllTasks(t *testing.T) {
	tests := []struct {
		name     string
		existing []model.Task
		expected []model.Task
		wantErr  bool
	}{
		{
			name: "success: get all existing tasks",
			existing: []model.Task{
				{
					ID:   "1",
					Name: "Task 1",
				},
				{
					ID:   "2",
					Name: "Task 2",
				},
			},
			expected: []model.Task{
				{
					ID:   "1",
					Name: "Task 1",
				},
				{
					ID:   "2",
					Name: "Task 2",
				},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := setupTestRepository(t)
			taskRepo := NewJSONTaskRepository(repo)

			for _, existingTask := range tt.existing {
				_, err := taskRepo.CreateTask(existingTask)
				if err != nil {
					t.Fatalf("Failed to create existing task: %v", err)
				}
			}

			got, err := taskRepo.GetAllTasks()
			if err != nil {
				if !tt.wantErr {
					t.Errorf("GetAllTasks() error = %v, wantErr %v", err, tt.wantErr)
				}
				return
			}

			assert.Equal(t, tt.expected, got)
		})
	}
}

func TestUpdateTask(t *testing.T) {
	tests := []struct {
		name         string
		existing     []model.Task
		taskToUpdate model.Task
		expected     model.Task
		wantErr      bool
	}{
		{
			name: "success: update existing task",
			existing: []model.Task{
				{
					ID:   "1",
					Name: "Task 1",
				},
			},
			taskToUpdate: model.Task{
				ID:   "1",
				Name: "Updated Task 1",
			},
			expected: model.Task{
				ID:   "1",
				Name: "Updated Task 1",
			},
			wantErr: false,
		},
		{
			name: "failure: task to update not found",
			existing: []model.Task{
				{
					ID:   "1",
					Name: "Task 1",
				},
			},
			taskToUpdate: model.Task{
				ID:   "non-existent-id",
				Name: "Non-existent Task",
			},
			expected: model.Task{},
			wantErr:  true,
		},
		{
			name: "failure: task name already exists",
			existing: []model.Task{
				{
					ID:   "1",
					Name: "Task 1",
				},
				{
					ID:   "2",
					Name: "Task 2",
				},
			},
			taskToUpdate: model.Task{
				ID:   "1",
				Name: "Task 2", // Name already exists for another task
			},
			expected: model.Task{},
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := setupTestRepository(t)
			taskRepo := NewJSONTaskRepository(repo)

			for _, existingTask := range tt.existing {
				_, err := taskRepo.CreateTask(existingTask)
				if err != nil {
					t.Fatalf("Failed to create existing task: %v", err)
				}
			}

			got, err := taskRepo.UpdateTask(tt.taskToUpdate)
			if err != nil {
				if !tt.wantErr {
					t.Errorf("UpdateTask() error = %v, wantErr %v", err, tt.wantErr)
				}
				return
			}
			assert.Equal(t, tt.expected, got)

			// verify updated
			newTask, err := taskRepo.GetTaskByID(tt.taskToUpdate.ID)
			if err != nil {
				t.Fatalf("Failed to get updated task: %v", err)
			}
			assert.Equal(t, tt.expected, newTask)
		})
	}
}

func TestDeleteTask(t *testing.T) {
	tests := []struct {
		name     string
		existing []model.Task
		id string
		wantErr  bool
	}{
		{
			name: "success: delete existing task",
			existing: []model.Task{
				{
					ID:   "1",
					Name: "Task 1",
				},
			},
			id: "1",
			wantErr: false,
		},
		{
			name: "failure: task to delete not found",
			existing: []model.Task{
				{
					ID:   "1",
					Name: "Task 1",
				},
				{
					ID:   "2",
					Name: "Task 2",
				},
				{
					ID:   "3",
					Name: "Task 3",
				},
			},
			id:       "non-existent-id",
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := setupTestRepository(t)
			taskRepo := NewJSONTaskRepository(repo)

			for _, existingTask := range tt.existing {
				_, err := taskRepo.CreateTask(existingTask)
				if err != nil {
					t.Fatalf("Failed to create existing task: %v", err)
				}
			}

			err := taskRepo.DeleteTask(tt.id)
			if err != nil {
				if !tt.wantErr {
					t.Errorf("DeleteTask() error = %v, wantErr %v", err, tt.wantErr)
				}
				return
			}

			// verify deleted
			_, err = taskRepo.GetTaskByID(tt.id)
			if err == nil {
				t.Errorf("Expected error when getting deleted task, got nil")
			}
		})
	}
}
