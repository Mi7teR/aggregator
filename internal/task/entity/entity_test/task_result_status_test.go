package entity_test

import (
	"reflect"
	"testing"

	"github.com/Mi7teR/aggregator/internal/task/entity"
)

func TestTaskResultStatus_MarshalJSON(t *testing.T) {
	tests := []struct {
		name    string
		t       entity.TaskResultStatus
		want    []byte
		wantErr bool
	}{
		{
			"marshall status new",
			entity.TaskStatusNew,
			[]byte(`"new"`),
			false,
		},
		{
			"marshall status in_process",
			entity.TaskStatusInProcess,
			[]byte(`"in_process"`),
			false,
		},
		{
			"marshall status error",
			entity.TaskStatusError,
			[]byte(`"error"`),
			false,
		},
		{
			"marshall status new",
			entity.TaskStatusDone,
			[]byte(`"done"`),
			false,
		},
		{
			"marshall invalid status error",
			0,
			nil,
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.t.MarshalJSON()
			if (err != nil) != tt.wantErr {
				t.Errorf("MarshalJSON() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("MarshalJSON() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTaskResultStatus_String(t *testing.T) {
	tests := []struct {
		name string
		t    entity.TaskResultStatus
		want string
	}{
		{
			"string from task status new",
			entity.TaskStatusNew,
			"new",
		},
		{
			"string from task status in_process",
			entity.TaskStatusInProcess,
			"in_process",
		},
		{
			"string from task status error",
			entity.TaskStatusError,
			"error",
		},
		{
			"string from task status done",
			entity.TaskStatusDone,
			"done",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.t.String(); got != tt.want {
				t.Errorf("String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTaskResultStatus_UnmarshalJSON(t *testing.T) {
	type args struct {
		i []byte
	}
	tests := []struct {
		name    string
		t       entity.TaskResultStatus
		args    args
		wantErr bool
	}{
		{
			"unmarshall task status new",
			entity.TaskStatusNew,
			args{
				i: []byte(`"new"`),
			},
			false,
		},
		{
			"unmarshall task status in_process",
			entity.TaskStatusInProcess,
			args{
				i: []byte(`"in_process"`),
			},
			false,
		},
		{
			"unmarshall task status error",
			entity.TaskStatusError,
			args{
				i: []byte(`"error"`),
			},
			false,
		},
		{
			"unmarshall task status done",
			entity.TaskStatusDone,
			args{
				i: []byte(`"done"`),
			},
			false,
		},
		{
			"unmarshall error prefix not found",
			0,
			args{
				i: []byte(`new"`),
			},
			true,
		},
		{
			"unmarshall error suffix not found",
			0,
			args{
				i: []byte(`"new`),
			},
			true,
		},
		{
			"unmarshall error invalid status",
			0,
			args{
				i: nil,
			},
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.t.UnmarshalJSON(tt.args.i); (err != nil) != tt.wantErr {
				t.Errorf("UnmarshalJSON() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
