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
)

type TaskType string

const (
	TaskTypeMangaGeneration TaskType = "manga_generation"
)

type Task struct {
	ID        string
	Type      TaskType
	Status    TaskStatus
	Input     map[string]interface{}
	Result    map[string]interface{}
	Error     string
	Progress  int
	Message   string
	CreatedAt time.Time
	UpdatedAt time.Time
	StartedAt *time.Time
	CompletedAt *time.Time
}

func NewTask(taskType TaskType, input map[string]interface{}) *Task {
	now := time.Now()
	return &Task{
		ID:        uuid.New().String(),
		Type:      taskType,
		Status:    TaskStatusPending,
		Input:     input,
		Progress:  0,
		CreatedAt: now,
		UpdatedAt: now,
	}
}

func (t *Task) Start() {
	now := time.Now()
	t.Status = TaskStatusProcessing
	t.StartedAt = &now
	t.UpdatedAt = now
}

func (t *Task) Complete(result map[string]interface{}) {
	now := time.Now()
	t.Status = TaskStatusCompleted
	t.Result = result
	t.Progress = 100
	t.CompletedAt = &now
	t.UpdatedAt = now
}

func (t *Task) Fail(err error) {
	now := time.Now()
	t.Status = TaskStatusFailed
	t.Error = err.Error()
	t.CompletedAt = &now
	t.UpdatedAt = now
}

func (t *Task) UpdateProgress(progress int, message string) {
	t.Progress = progress
	t.Message = message
	t.UpdatedAt = time.Now()
}
