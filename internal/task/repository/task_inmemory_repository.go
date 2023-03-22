package repository

import (
	"context"
	"errors"
	"sync"

	"github.com/Mi7teR/aggregator/internal/task/entity"
	"github.com/google/uuid"
)

type TaskInMemoryRepository struct {
	mu   sync.RWMutex
	data map[string]entity.TaskResult
}

var ErrNotFound = errors.New("task result not found")

func NewTaskInMemoryRepository() *TaskInMemoryRepository {
	return &TaskInMemoryRepository{data: make(map[string]entity.TaskResult)}
}

func (t *TaskInMemoryRepository) Create(ctx context.Context, task *entity.Task) (string, error) {
	t.mu.Lock()
	defer t.mu.Unlock()

	newTask := entity.TaskResult{
		ID:             uuid.New().String(),
		Status:         entity.TaskStatusNew,
		HTTPStatusCode: 0,
		Headers:        nil,
		Length:         0,
	}

	t.data[newTask.ID] = newTask

	return newTask.ID, nil
}

func (t *TaskInMemoryRepository) GetByID(ctx context.Context, id string) (*entity.TaskResult, error) {
	t.mu.RLock()
	defer t.mu.RUnlock()

	v, ok := t.data[id]
	if !ok {
		return nil, ErrNotFound
	}

	return &v, nil
}

func (t *TaskInMemoryRepository) Update(ctx context.Context, res *entity.TaskResult) error {
	t.mu.Lock()
	defer t.mu.Unlock()

	_, ok := t.data[res.ID]
	if !ok {
		return ErrNotFound
	}

	t.data[res.ID] = *res

	return nil
}
