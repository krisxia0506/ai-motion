package task

import "context"

// Repository Task 仓储接口
type Repository interface {
	// Save 保存任务（创建或更新）
	Save(ctx context.Context, task *Task) error

	// FindByID 根据ID查找任务
	FindByID(ctx context.Context, taskID string) (*Task, error)

	// FindByIDAndUserID 根据ID和用户ID查找任务（用于权限验证）
	FindByIDAndUserID(ctx context.Context, taskID, userID string) (*Task, error)

	// FindByUserID 查询用户的所有任务（分页）
	FindByUserID(ctx context.Context, userID string, page, pageSize int, status string) ([]*Task, int, error)

	// Delete 删除任务
	Delete(ctx context.Context, taskID string) error
}
