package memory

import (
	"context"
	"errors"
	"sync"

	"github.com/xiajiayi/ai-motion/internal/domain/task"
)

var (
	ErrTaskNotFound = errors.New("task not found")
)

type TaskRepository struct {
	tasks map[string]*task.Task
	mu    sync.RWMutex
}

func NewTaskRepository() task.TaskRepository {
	return &TaskRepository{
		tasks: make(map[string]*task.Task),
	}
}

func (r *TaskRepository) Save(ctx context.Context, t *task.Task) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.tasks[t.ID] = t
	return nil
}

func (r *TaskRepository) FindByID(ctx context.Context, id string) (*task.Task, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	t, exists := r.tasks[id]
	if !exists {
		return nil, ErrTaskNotFound
	}

	return t, nil
}

func (r *TaskRepository) Update(ctx context.Context, t *task.Task) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.tasks[t.ID]; !exists {
		return ErrTaskNotFound
	}

	r.tasks[t.ID] = t
	return nil
}

func (r *TaskRepository) ListByStatus(ctx context.Context, status task.TaskStatus) ([]*task.Task, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var tasks []*task.Task
	for _, t := range r.tasks {
		if t.Status == status {
			tasks = append(tasks, t)
		}
	}

	return tasks, nil
}
