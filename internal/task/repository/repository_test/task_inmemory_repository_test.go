package repository_test

import (
	"context"
	"net/http"
	"reflect"
	"testing"

	"github.com/Mi7teR/aggregator/internal/task/entity"
	"github.com/Mi7teR/aggregator/internal/task/repository"
	"github.com/google/uuid"
)

func TestNewTaskInMemoryRepository(t *testing.T) {
	tests := []struct {
		name string
		want *repository.TaskInMemoryRepository
	}{
		{
			name: "repository created",
			want: repository.NewTaskInMemoryRepository(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := repository.NewTaskInMemoryRepository(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewTaskInMemoryRepository() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTaskInMemoryRepository_Create(t *testing.T) {
	type fields struct {
		repo *repository.TaskInMemoryRepository
	}
	type args struct {
		ctx  context.Context
		task *entity.Task
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name:   "task created",
			fields: fields{repo: repository.NewTaskInMemoryRepository()},
			args: args{
				ctx:  context.Background(),
				task: &entity.Task{},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.fields.repo.Create(tt.args.ctx, tt.args.task)
			if (err != nil) != tt.wantErr {
				t.Errorf("Create() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if _, err = uuid.Parse(got); err != nil {
				t.Errorf("created invalid uuid")
			}
		})
	}
}

func TestTaskInMemoryRepository_GetByID(t *testing.T) {
	repo := repository.NewTaskInMemoryRepository()
	taskID, _ := repo.Create(context.Background(), &entity.Task{})

	type fields struct {
		repo *repository.TaskInMemoryRepository
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
			name: "get existent task result",
			fields: fields{
				repo: repo,
			},
			args: args{
				ctx: context.Background(),
				id:  taskID,
			},
			want: &entity.TaskResult{
				ID:             taskID,
				Status:         entity.TaskStatusNew,
				HTTPStatusCode: 0,
				Headers:        nil,
				Length:         0,
			},
			wantErr: false,
		},
		{
			name: "get non-existent task result",
			fields: fields{
				repo: repo,
			},
			args: args{
				ctx: context.Background(),
				id:  "invalid task id",
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.fields.repo.GetByID(tt.args.ctx, tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetByID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetByID() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTaskInMemoryRepository_Update(t *testing.T) {
	repo := repository.NewTaskInMemoryRepository()
	taskID, _ := repo.Create(context.Background(), &entity.Task{})

	type fields struct {
		repo *repository.TaskInMemoryRepository
	}
	type args struct {
		ctx context.Context
		res *entity.TaskResult
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "update existent task result",
			fields: fields{
				repo: repo,
			},
			args: args{
				ctx: context.Background(),
				res: &entity.TaskResult{
					ID:             taskID,
					Status:         entity.TaskStatusDone,
					HTTPStatusCode: http.StatusOK,
					Headers:        nil,
					Length:         20,
				},
			},
			wantErr: false,
		},
		{
			name: "update non-existent task result",
			fields: fields{
				repo: repo,
			},
			args: args{
				ctx: context.Background(),
				res: &entity.TaskResult{
					ID:             "invalid-id",
					Status:         entity.TaskStatusDone,
					HTTPStatusCode: http.StatusOK,
					Headers:        nil,
					Length:         20,
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.fields.repo.Update(tt.args.ctx, tt.args.res); (err != nil) != tt.wantErr {
				t.Errorf("Update() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
