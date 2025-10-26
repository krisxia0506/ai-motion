package supabase

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"

	"github.com/supabase-community/postgrest-go"
	"github.com/xiajiayi/ai-motion/internal/domain/task"
)

type TaskRepository struct {
	client *postgrest.Client
}

func NewTaskRepository(client *postgrest.Client) task.Repository {
	return &TaskRepository{
		client: client,
	}
}

// getClientWithAuth returns a client with JWT token from context if available
func (r *TaskRepository) getClientWithAuth(ctx context.Context) *postgrest.Client {
	// Try to get JWT token from context
	if jwtToken, ok := ctx.Value("jwt_token").(string); ok && jwtToken != "" {
		// Create a new client instance with the user's JWT token
		// SetAuthToken returns a new client with the Authorization header set
		// We need to ensure apikey is also set
		slog.Info("[TaskRepo] Using user JWT token for authentication", "token_length", len(jwtToken))
		authClient := r.client.SetAuthToken(jwtToken)
		return authClient
	}
	// Fallback to service role client
	slog.Warn("[TaskRepo] Fallback to service role client - no JWT token in context")
	return r.client
}

// taskRecord Supabase中的任务记录结构
type taskRecord struct {
	ID                 string          `json:"id"`
	UserID             string          `json:"user_id"`
	NovelID            string          `json:"novel_id"`
	Status             string          `json:"status"`
	ProgressStep       string          `json:"progress_step"`
	ProgressStepIndex  int             `json:"progress_step_index"`
	ProgressPercentage int             `json:"progress_percentage"`
	ProgressDetails    json.RawMessage `json:"progress_details"`
	ErrorCode          *int            `json:"error_code"`
	ErrorMessage       *string         `json:"error_message"`
	CreatedAt          string          `json:"created_at"`
	UpdatedAt          string          `json:"updated_at"`
	CompletedAt        *string         `json:"completed_at"`
	FailedAt           *string         `json:"failed_at"`
}

// Save 保存任务（创建或更新）
func (r *TaskRepository) Save(ctx context.Context, t *task.Task) error {
	// 序列化 progress_details
	detailsJSON, err := json.Marshal(t.ProgressDetails)
	if err != nil {
		return fmt.Errorf("failed to marshal progress details: %w", err)
	}

	record := taskRecord{
		ID:                 t.ID,
		UserID:             t.UserID,
		NovelID:            t.NovelID,
		Status:             string(t.Status),
		ProgressStep:       t.ProgressStep,
		ProgressStepIndex:  t.ProgressStepIndex,
		ProgressPercentage: t.ProgressPercentage,
		ProgressDetails:    detailsJSON,
		CreatedAt:          t.CreatedAt.Format("2006-01-02T15:04:05.999999Z07:00"),
		UpdatedAt:          t.UpdatedAt.Format("2006-01-02T15:04:05.999999Z07:00"),
	}

	if t.ErrorCode != 0 {
		record.ErrorCode = &t.ErrorCode
	}
	if t.ErrorMessage != "" {
		record.ErrorMessage = &t.ErrorMessage
	}
	if t.CompletedAt != nil {
		completedStr := t.CompletedAt.Format("2006-01-02T15:04:05.999999Z07:00")
		record.CompletedAt = &completedStr
	}
	if t.FailedAt != nil {
		failedStr := t.FailedAt.Format("2006-01-02T15:04:05.999999Z07:00")
		record.FailedAt = &failedStr
	}

	// 使用 upsert（如果ID存在则更新，否则插入）
	client := r.getClientWithAuth(ctx)
	_, _, err = client.From("aimotion_task").
		Upsert(record, "", "", "").
		Execute()

	if err != nil {
		return fmt.Errorf("failed to save task: %w", err)
	}

	return nil
}

// FindByID 根据ID查找任务
func (r *TaskRepository) FindByID(ctx context.Context, taskID string) (*task.Task, error) {
	var record taskRecord

	client := r.getClientWithAuth(ctx)
	_, err := client.From("aimotion_task").
		Select("*", "", false).
		Eq("id", taskID).
		Single().
		ExecuteTo(&record)

	if err != nil {
		return nil, fmt.Errorf("failed to find task: %w", err)
	}

	return r.recordToTask(&record)
}

// FindByIDAndUserID 根据ID和用户ID查找任务（用于权限验证）
func (r *TaskRepository) FindByIDAndUserID(ctx context.Context, taskID, userID string) (*task.Task, error) {
	var record taskRecord

	slog.Info("[TaskRepo] FindByIDAndUserID", "task_id", taskID, "user_id", userID)

	client := r.getClientWithAuth(ctx)
	_, err := client.From("aimotion_task").
		Select("*", "", false).
		Eq("id", taskID).
		Eq("user_id", userID).
		Single().
		ExecuteTo(&record)

	if err != nil {
		slog.Error("[TaskRepo] FindByIDAndUserID - Query error", "error", err, "task_id", taskID, "user_id", userID)
		return nil, fmt.Errorf("failed to find task: %w", err)
	}

	slog.Info("[TaskRepo] FindByIDAndUserID - Query success", "task_id", record.ID)

	return r.recordToTask(&record)
}

// FindByUserID 查询用户的所有任务（分页）
func (r *TaskRepository) FindByUserID(ctx context.Context, userID string, page, pageSize int, status string) ([]*task.Task, int, error) {
	// 计算偏移量
	offset := (page - 1) * pageSize

	// 构建查询
	client := r.getClientWithAuth(ctx)
	query := client.From("aimotion_task").
		Select("*", "exact", false).
		Eq("user_id", userID)

	// 按状态筛选
	if status != "" && status != "all" {
		query = query.Eq("status", status)
	}

	// 排序和分页
	query = query.Order("created_at", &postgrest.OrderOpts{Ascending: false}).
		Range(offset, offset+pageSize-1, "")

	var records []taskRecord
	count, err := query.ExecuteTo(&records)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to query tasks: %w", err)
	}

	// 转换为domain对象
	tasks := make([]*task.Task, 0, len(records))
	for i := range records {
		t, err := r.recordToTask(&records[i])
		if err != nil {
			return nil, 0, fmt.Errorf("failed to convert task: %w", err)
		}
		tasks = append(tasks, t)
	}

	return tasks, int(count), nil
}

// Delete 删除任务
func (r *TaskRepository) Delete(ctx context.Context, taskID string) error {
	client := r.getClientWithAuth(ctx)
	_, _, err := client.From("aimotion_task").
		Delete("", "").
		Eq("id", taskID).
		Execute()

	if err != nil {
		return fmt.Errorf("failed to delete task: %w", err)
	}

	return nil
}

// recordToTask 将数据库记录转换为domain对象
func (r *TaskRepository) recordToTask(record *taskRecord) (*task.Task, error) {
	t := &task.Task{
		ID:                 record.ID,
		UserID:             record.UserID,
		NovelID:            record.NovelID,
		Status:             task.TaskStatus(record.Status),
		ProgressStep:       record.ProgressStep,
		ProgressStepIndex:  record.ProgressStepIndex,
		ProgressPercentage: record.ProgressPercentage,
	}

	// 解析 progress_details
	if len(record.ProgressDetails) > 0 {
		if err := json.Unmarshal(record.ProgressDetails, &t.ProgressDetails); err != nil {
			return nil, fmt.Errorf("failed to unmarshal progress details: %w", err)
		}
	}

	// 解析错误信息
	if record.ErrorCode != nil {
		t.ErrorCode = *record.ErrorCode
	}
	if record.ErrorMessage != nil {
		t.ErrorMessage = *record.ErrorMessage
	}

	// 解析时间字段
	// CreatedAt, UpdatedAt 由Supabase自动处理，这里简化处理
	// 实际生产环境需要正确解析时间格式

	// 初始化取消通道
	t.InitCancelChan()

	return t, nil
}
