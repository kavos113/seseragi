package usecase

//go:generate mockgen -source=workflow.go -destination=./mock_usecase/workflow_mock.go -package=mock_usecase
//go:generate mockgen -source=task.go -destination=./mock_usecase/task_mock.go -package=mock_usecase
//go:generate mockgen -source=workflow_run.go -destination=./mock_usecase/workflow_run_mock.go -package=mock_usecase
