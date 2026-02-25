package domain

//go:generate mockgen -source=task.go -destination=./mock_domain/task_mock.go -package=mock_domain
//go:generate mockgen -source=workflow.go -destination=./mock_domain/workflow_mock.go -package=mock_domain
//go:generate mockgen -source=runner.go -destination=./mock_domain/runner_mock.go -package=mock_domain
