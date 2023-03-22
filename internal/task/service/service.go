package service

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/Mi7teR/aggregator/internal/task/entity"
)

type Service struct {
	repo    Repository
	timeout time.Duration
}

func NewService(repo Repository, timeout time.Duration) *Service {
	return &Service{repo: repo, timeout: timeout}
}

func (s *Service) GetTaskResult(ctx context.Context, id string) (*entity.TaskResult, error) {
	res, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (s *Service) AddTask(ctx context.Context, task *entity.Task) (string, error) {
	taskID, err := s.repo.Create(ctx, task)
	if err != nil {
		return "", err
	}

	go s.Execute(taskID, task)

	return taskID, nil
}

func (s *Service) Execute(id string, task *entity.Task) {
	ctx, cancel := context.WithTimeout(context.Background(), s.timeout)
	defer cancel()

	err := s.repo.Update(ctx, &entity.TaskResult{
		ID:     id,
		Status: entity.TaskStatusInProcess,
	})
	if err != nil {
		log.Println(err)
		return
	}

	req, err := http.NewRequestWithContext(ctx, task.Method.String(), task.URL, nil)
	if err != nil {
		errUpdate := s.repo.Update(ctx, &entity.TaskResult{
			ID:     id,
			Status: entity.TaskStatusError,
		})

		if errUpdate != nil {
			log.Println(err)
			return
		}
	}

	for i := range task.Headers {
		req.Header.Add(i, task.Headers[i])
	}

	client := &http.Client{}

	res, err := client.Do(req)
	if err != nil {
		errUpdate := s.repo.Update(ctx, &entity.TaskResult{
			ID:     id,
			Status: entity.TaskStatusError,
		})

		if errUpdate != nil {
			log.Println(err)
		}
		return
	}

	defer res.Body.Close()

	err = s.repo.Update(ctx, &entity.TaskResult{
		ID:             id,
		Status:         entity.TaskStatusDone,
		HTTPStatusCode: res.StatusCode,
		Headers:        res.Header,
		Length:         res.ContentLength,
	})
	if err != nil {
		log.Println(err)
		return
	}
}
