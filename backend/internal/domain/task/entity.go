package task

import (
	"time"

	"github.com/google/uuid"
)

type TaskStatus string

const (
	TaskStatusPending    TaskStatus = "pending"
	TaskStatusProcessing TaskStatus = "processing"
	TaskStatusCompleted  TaskStatus = "completed"
	TaskStatusFailed     TaskStatus = "failed"
	TaskStatusCancelled  TaskStatus = "cancelled"
)

type ProgressDetails struct {
	CharactersExtracted int `json:"characters_extracted"`
	CharactersGenerated int `json:"characters_generated"`
	ScenesDivided       int `json:"scenes_divided"`
	ScenesGenerated     int `json:"scenes_generated"`
}

type Task struct {
	ID                 string          `json:"id"`
	UserID             string          `json:"user_id"` // 用户ID，用于数据隔离
	NovelID            string          `json:"novel_id"`
	Status             TaskStatus      `json:"status"`
	ProgressStep       string          `json:"progress_step"`
	ProgressStepIndex  int             `json:"progress_step_index"`
	ProgressPercentage int             `json:"progress_percentage"`
	ProgressDetails    ProgressDetails `json:"progress_details"`
	ErrorCode          int             `json:"error_code,omitempty"`
	ErrorMessage       string          `json:"error_message,omitempty"`
	CreatedAt          time.Time       `json:"created_at"`
	UpdatedAt          time.Time       `json:"updated_at"`
	CompletedAt        *time.Time      `json:"completed_at,omitempty"`
	FailedAt           *time.Time      `json:"failed_at,omitempty"`
	cancelChan         chan struct{}   `json:"-"` // 取消信号通道
}

// NewTask 创建新任务
func NewTask(userID, novelID string) *Task {
	return &Task{
		ID:                 uuid.New().String(),
		UserID:             userID,
		NovelID:            novelID,
		Status:             TaskStatusPending,
		ProgressStep:       "等待中",
		ProgressStepIndex:  0,
		ProgressPercentage: 0,
		ProgressDetails:    ProgressDetails{},
		CreatedAt:          time.Now(),
		UpdatedAt:          time.Now(),
		cancelChan:         make(chan struct{}),
	}
}

// UpdateProgress 更新任务进度
func (t *Task) UpdateProgress(step string, stepIndex int, percentage int, details ProgressDetails) {
	t.Status = TaskStatusProcessing
	t.ProgressStep = step
	t.ProgressStepIndex = stepIndex
	t.ProgressPercentage = percentage
	t.ProgressDetails = details
	t.UpdatedAt = time.Now()
}

// MarkCompleted 标记任务完成
func (t *Task) MarkCompleted() {
	t.Status = TaskStatusCompleted
	t.ProgressStep = "完成"
	t.ProgressStepIndex = 6
	t.ProgressPercentage = 100
	now := time.Now()
	t.CompletedAt = &now
	t.UpdatedAt = now
}

// MarkFailed 标记任务失败
func (t *Task) MarkFailed(errCode int, errMsg string) {
	t.Status = TaskStatusFailed
	t.ErrorCode = errCode
	t.ErrorMessage = errMsg
	now := time.Now()
	t.FailedAt = &now
	t.UpdatedAt = now
}

// Cancel 取消任务
func (t *Task) Cancel() {
	t.Status = TaskStatusCancelled
	t.UpdatedAt = time.Now()
	if t.cancelChan != nil {
		close(t.cancelChan)
	}
}

// IsCancelled 检查任务是否已取消
func (t *Task) IsCancelled() bool {
	if t.cancelChan == nil {
		return false
	}
	select {
	case <-t.cancelChan:
		return true
	default:
		return false
	}
}

// InitCancelChan 初始化取消通道（用于从数据库加载后）
func (t *Task) InitCancelChan() {
	if t.cancelChan == nil {
		t.cancelChan = make(chan struct{})
	}
}

// IsRetryable 判断任务失败后是否可以重试
func (t *Task) IsRetryable() bool {
	if t.Status != TaskStatusFailed {
		return false
	}
	// 40001 是 AI 服务调用失败，通常可以重试
	return t.ErrorCode >= 40001 && t.ErrorCode < 50000
}
