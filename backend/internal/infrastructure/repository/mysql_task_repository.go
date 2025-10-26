package repository

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/xiajiayi/ai-motion/backend/internal/domain/task"
)

type MySQLTaskRepository struct {
	db *sql.DB
}

func NewMySQLTaskRepository(db *sql.DB) *MySQLTaskRepository {
	return &MySQLTaskRepository{
		db: db,
	}
}

// Save 保存任务（创建或更新）
func (r *MySQLTaskRepository) Save(ctx context.Context, t *task.Task) error {
	// 序列化 progress_details 为 JSON
	detailsJSON, err := json.Marshal(t.ProgressDetails)
	if err != nil {
		return fmt.Errorf("failed to marshal progress details: %w", err)
	}

	query := `
		INSERT INTO tasks (
			id, user_id, novel_id, status, progress_step, progress_step_index,
			progress_percentage, progress_details, error_code, error_message,
			created_at, updated_at, completed_at, failed_at
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
		ON DUPLICATE KEY UPDATE
			status = VALUES(status),
			progress_step = VALUES(progress_step),
			progress_step_index = VALUES(progress_step_index),
			progress_percentage = VALUES(progress_percentage),
			progress_details = VALUES(progress_details),
			error_code = VALUES(error_code),
			error_message = VALUES(error_message),
			updated_at = VALUES(updated_at),
			completed_at = VALUES(completed_at),
			failed_at = VALUES(failed_at)
	`

	_, err = r.db.ExecContext(
		ctx,
		query,
		t.ID, t.UserID, t.NovelID, t.Status, t.ProgressStep, t.ProgressStepIndex,
		t.ProgressPercentage, detailsJSON, t.ErrorCode, t.ErrorMessage,
		t.CreatedAt, t.UpdatedAt, t.CompletedAt, t.FailedAt,
	)

	if err != nil {
		return fmt.Errorf("failed to save task: %w", err)
	}

	return nil
}

// FindByID 根据ID查找任务
func (r *MySQLTaskRepository) FindByID(ctx context.Context, taskID string) (*task.Task, error) {
	query := `
		SELECT id, user_id, novel_id, status, progress_step, progress_step_index,
			progress_percentage, progress_details, error_code, error_message,
			created_at, updated_at, completed_at, failed_at
		FROM tasks
		WHERE id = ?
	`

	row := r.db.QueryRowContext(ctx, query, taskID)
	return r.scanTask(row)
}

// FindByIDAndUserID 根据ID和用户ID查找任务（用于权限验证）
func (r *MySQLTaskRepository) FindByIDAndUserID(ctx context.Context, taskID, userID string) (*task.Task, error) {
	query := `
		SELECT id, user_id, novel_id, status, progress_step, progress_step_index,
			progress_percentage, progress_details, error_code, error_message,
			created_at, updated_at, completed_at, failed_at
		FROM tasks
		WHERE id = ? AND user_id = ?
	`

	row := r.db.QueryRowContext(ctx, query, taskID, userID)
	return r.scanTask(row)
}

// FindByUserID 查询用户的所有任务（分页）
func (r *MySQLTaskRepository) FindByUserID(ctx context.Context, userID string, page, pageSize int, status string) ([]*task.Task, int, error) {
	// 计算偏移量
	offset := (page - 1) * pageSize

	// 构建查询条件
	queryCondition := "WHERE user_id = ?"
	args := []interface{}{userID}

	if status != "" && status != "all" {
		queryCondition += " AND status = ?"
		args = append(args, status)
	}

	// 查询总数
	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM tasks %s", queryCondition)
	var total int
	err := r.db.QueryRowContext(ctx, countQuery, args...).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count tasks: %w", err)
	}

	// 查询任务列表
	query := fmt.Sprintf(`
		SELECT id, user_id, novel_id, status, progress_step, progress_step_index,
			progress_percentage, progress_details, error_code, error_message,
			created_at, updated_at, completed_at, failed_at
		FROM tasks
		%s
		ORDER BY created_at DESC
		LIMIT ? OFFSET ?
	`, queryCondition)

	args = append(args, pageSize, offset)

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to query tasks: %w", err)
	}
	defer rows.Close()

	var tasks []*task.Task
	for rows.Next() {
		t, err := r.scanTaskFromRows(rows)
		if err != nil {
			return nil, 0, err
		}
		tasks = append(tasks, t)
	}

	if err = rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("error iterating tasks: %w", err)
	}

	return tasks, total, nil
}

// Delete 删除任务
func (r *MySQLTaskRepository) Delete(ctx context.Context, taskID string) error {
	query := "DELETE FROM tasks WHERE id = ?"
	result, err := r.db.ExecContext(ctx, query, taskID)
	if err != nil {
		return fmt.Errorf("failed to delete task: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return errors.New("task not found")
	}

	return nil
}

// scanTask 从单行扫描任务
func (r *MySQLTaskRepository) scanTask(row *sql.Row) (*task.Task, error) {
	var t task.Task
	var detailsJSON []byte

	err := row.Scan(
		&t.ID, &t.UserID, &t.NovelID, &t.Status, &t.ProgressStep, &t.ProgressStepIndex,
		&t.ProgressPercentage, &detailsJSON, &t.ErrorCode, &t.ErrorMessage,
		&t.CreatedAt, &t.UpdatedAt, &t.CompletedAt, &t.FailedAt,
	)

	if err == sql.ErrNoRows {
		return nil, errors.New("task not found")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to scan task: %w", err)
	}

	// 反序列化 progress_details
	if len(detailsJSON) > 0 {
		if err := json.Unmarshal(detailsJSON, &t.ProgressDetails); err != nil {
			return nil, fmt.Errorf("failed to unmarshal progress details: %w", err)
		}
	}

	// 初始化取消通道
	t.InitCancelChan()

	return &t, nil
}

// scanTaskFromRows 从多行扫描任务
func (r *MySQLTaskRepository) scanTaskFromRows(rows *sql.Rows) (*task.Task, error) {
	var t task.Task
	var detailsJSON []byte

	err := rows.Scan(
		&t.ID, &t.UserID, &t.NovelID, &t.Status, &t.ProgressStep, &t.ProgressStepIndex,
		&t.ProgressPercentage, &detailsJSON, &t.ErrorCode, &t.ErrorMessage,
		&t.CreatedAt, &t.UpdatedAt, &t.CompletedAt, &t.FailedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to scan task: %w", err)
	}

	// 反序列化 progress_details
	if len(detailsJSON) > 0 {
		if err := json.Unmarshal(detailsJSON, &t.ProgressDetails); err != nil {
			return nil, fmt.Errorf("failed to unmarshal progress details: %w", err)
		}
	}

	// 初始化取消通道
	t.InitCancelChan()

	return &t, nil
}
