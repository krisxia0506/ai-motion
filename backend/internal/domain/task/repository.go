package task

import "context"

type TaskRepository interface {
	Save(ctx context.Context, task *Task) error
	FindByID(ctx context.Context, id string) (*Task, error)
	Update(ctx context.Context, task *Task) error
	ListByStatus(ctx context.Context, status TaskStatus) ([]*Task, error)
}
