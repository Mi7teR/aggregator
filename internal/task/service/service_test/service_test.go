package service_test

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"sync"
	"testing"
	"time"

	"github.com/Mi7teR/aggregator/internal/task/entity"
	"github.com/Mi7teR/aggregator/internal/task/repository"
	"github.com/Mi7teR/aggregator/internal/task/service"
	"github.com/google/uuid"
)

func TestNewService(t *testing.T) {
	type args struct {
		repo    service.Repository
		timeout time.Duration
	}
	tests := []struct {
		name string
		args args
		want *service.Service
	}{
		{
			"create new service",
			args{
				repo:    repository.NewTaskInMemoryRepository(),
				timeout: time.Second * 30,
			},
			service.NewService(repository.NewTaskInMemoryRepository(), time.Second*30),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := service.NewService(tt.args.repo, tt.args.timeout); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewService() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestService_AddTask(t *testing.T) {
	wg := &sync.WaitGroup{}
	wg.Add(1)
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer wg.Done()
		if r.Method != http.MethodPost {
			t.Errorf("Expected to request with method POST, got: %s", r.Method)
		}
		if r.URL.Path != "/test-path" {
			t.Errorf("Expected to request '/test-path', got: %s", r.URL.Path)
		}
		if r.Header.Get("Accept") != "application/json" {
			t.Errorf("Expected Accept: application/json header, got: %s", r.Header.Get("Accept"))
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`1234567890`)) //nolint:errcheck // we dont test it :)
	}))
	defer server.Close()

	task := &entity.Task{
		Method: entity.MethodPost,
		URL:    fmt.Sprintf("%s/test-path", server.URL),
		Headers: map[string]string{
			"Accept": "application/json",
		},
	}

	repo := repository.NewTaskInMemoryRepository()
	s := service.NewService(repo, time.Second*30)

	id, err := s.AddTask(context.Background(), task)
	if err != nil {
		t.Errorf("Expected to add task without errors, got: %s", err)
	}

	wg.Wait() // waiting for request end in goroutine
	if _, err = uuid.Parse(id); err != nil {
		t.Errorf("Expected to parse task id without errors, got: %s", err)
	}
}

func TestService_GetTaskResult(t *testing.T) {
	repo := repository.NewTaskInMemoryRepository()
	timeout := time.Second * 1000

	id, err := repo.Create(context.Background(), &entity.Task{
		Method:  entity.MethodGet,
		URL:     "https://google.com",
		Headers: nil,
	})
	if err != nil {
		t.Errorf("Expected to create new task result, got %s", err)
	}

	if _, err = uuid.Parse(id); err != nil {
		t.Errorf("Expected to parse task id without errors, got: %s", err)
	}

	type fields struct {
		repo    service.Repository
		timeout time.Duration
	}
	type args struct {
		ctx context.Context
		id  string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *entity.TaskResult
		wantErr bool
	}{
		{
			"get existent task",
			fields{
				repo:    repo,
				timeout: timeout,
			},
			args{
				ctx: context.Background(),
				id:  id,
			},
			&entity.TaskResult{
				ID:             id,
				Status:         entity.TaskStatusNew,
				HTTPStatusCode: 0,
				Headers:        nil,
				Length:         0,
			},
			false,
		},
		{
			"get non-existent task",
			fields{
				repo:    repo,
				timeout: timeout,
			},
			args{
				ctx: context.Background(),
				id:  "invalid-id",
			},
			nil,
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := service.NewService(tt.fields.repo, tt.fields.timeout)
			got, err := s.GetTaskResult(tt.args.ctx, tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetTaskResult() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetTaskResult() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestService_Execute(t *testing.T) {
	repo := repository.NewTaskInMemoryRepository()
	timeout := time.Second * 1000

	wg := &sync.WaitGroup{}
	wg.Add(1)
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer wg.Done()
		if r.Method != http.MethodPost {
			t.Errorf("Expected to request with method POST, got: %s", r.Method)
		}
		if r.URL.Path != "/test-path" {
			t.Errorf("Expected to request '/test-path', got: %s", r.URL.Path)
		}
		if r.Header.Get("Accept") != "application/json" {
			t.Errorf("Expected Accept: application/json header, got: %s", r.Header.Get("Accept"))
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`1234567890`)) //nolint:errcheck // we dont test it :)
	}))
	defer server.Close()

	task := &entity.Task{
		Method: entity.MethodPost,
		URL:    fmt.Sprintf("%s/test-path", server.URL),
		Headers: map[string]string{
			"Accept": "application/json",
		},
	}

	id, err := repo.Create(context.Background(), task)
	if err != nil {
		t.Errorf("Expected to create new task result, got %s", err)
	}

	if _, err = uuid.Parse(id); err != nil {
		t.Errorf("Expected to parse task id without errors, got: %s", err)
	}

	s := service.NewService(repo, timeout)

	s.Execute(id, task)

	wg.Wait()

	res, err := repo.GetByID(context.Background(), id)
	if err != nil {
		t.Errorf("expected to get task result, got %s", err)
	}

	taskResult := &entity.TaskResult{
		ID:             id,
		Status:         entity.TaskStatusDone,
		HTTPStatusCode: http.StatusOK,
		Headers: map[string][]string{
			"Content-Length": {
				"10",
			},
			"Content-Type": {
				"text/plain; charset=utf-8",
			},
			"Date": {
				time.Now().UTC().Format(http.TimeFormat),
			},
		},
		Length: 10,
	}
	if !reflect.DeepEqual(res, taskResult) {
		t.Errorf("Execute() got = %v, want %v", res, taskResult)
	}
}
