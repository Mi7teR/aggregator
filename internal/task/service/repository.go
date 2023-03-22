package service

import (
	"context"

	"github.com/Mi7teR/aggregator/internal/task/entity"
)

type Repository interface {
	Create(ctx context.Context, task *entity.Task) (string, error)
	GetByID(ctx context.Context, id string) (*entity.TaskResult, error)
	Update(ctx context.Context, res *entity.TaskResult) error
}
