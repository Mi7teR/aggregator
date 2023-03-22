package entity_test

import (
	"net/http"
	"reflect"
	"testing"

	"github.com/Mi7teR/aggregator/internal/task/entity"
)

func TestTaskMethod_MarshalJSON(t *testing.T) {
	tests := []struct {
		name    string
		t       entity.TaskMethod
		want    []byte
		wantErr bool
	}{
		{
			"marshal method GET",
			entity.MethodGet,
			[]byte(`"GET"`),
			false,
		},
		{
			"marshal method HEAD",
			entity.MethodHead,
			[]byte(`"HEAD"`),
			false,
		},
		{
			"marshal method POST",
			entity.MethodPost,
			[]byte(`"POST"`),
			false,
		},
		{
			"marshal method PUT",
			entity.MethodPut,
			[]byte(`"PUT"`),
			false,
		},
		{
			"marshal method PATCH",
			entity.MethodPatch,
			[]byte(`"PATCH"`),
			false,
		},
		{
			"marshal method DELETE",
			entity.MethodDelete,
			[]byte(`"DELETE"`),
			false,
		},
		{
			"marshal method CONNECT",
			entity.MethodConnect,
			[]byte(`"CONNECT"`),
			false,
		},
		{
			"marshal method OPTIONS",
			entity.MethodOptions,
			[]byte(`"OPTIONS"`),
			false,
		},
		{
			"marshal method TRACE",
			entity.MethodTrace,
			[]byte(`"TRACE"`),
			false,
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

func TestTaskMethod_String(t *testing.T) {
	tests := []struct {
		name string
		t    entity.TaskMethod
		want string
	}{
		{
			"MethodGet String()",
			entity.MethodGet,
			http.MethodGet,
		},
		{
			"MethodHead String()",
			entity.MethodHead,
			http.MethodHead,
		},
		{
			"MethodPost String()",
			entity.MethodPost,
			http.MethodPost,
		},
		{
			"MethodPut String()",
			entity.MethodPut,
			http.MethodPut,
		},
		{
			"MethodPatch String()",
			entity.MethodPatch,
			http.MethodPatch,
		},
		{
			"MethodDelete String()",
			entity.MethodDelete,
			http.MethodDelete,
		},
		{
			"MethodConnect String()",
			entity.MethodConnect,
			http.MethodConnect,
		},
		{
			"MethodOptions String()",
			entity.MethodOptions,
			http.MethodOptions,
		},
		{
			"MethodTrace String()",
			entity.MethodTrace,
			http.MethodTrace,
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

func TestTaskMethod_UnmarshalJSON(t *testing.T) {
	type args struct {
		i []byte
	}
	tests := []struct {
		name    string
		t       entity.TaskMethod
		args    args
		wantErr bool
	}{
		{
			"MethodGet UnmarshalJSON()",
			entity.MethodGet,
			args{
				i: []byte(`"GET"`),
			},
			false,
		},
		{
			"MethodHead UnmarshalJSON()",
			entity.MethodHead,
			args{
				i: []byte(`"HEAD"`),
			},
			false,
		},
		{
			"MethodPost UnmarshalJSON()",
			entity.MethodPost,
			args{
				i: []byte(`"POST"`),
			},
			false,
		},
		{
			"MethodPut UnmarshalJSON()",
			entity.MethodPut,
			args{
				i: []byte(`"PUT"`),
			},
			false,
		},
		{
			"MethodPatch UnmarshalJSON()",
			entity.MethodPatch,
			args{
				i: []byte(`"PATCH"`),
			},
			false,
		},
		{
			"MethodDelete UnmarshalJSON()",
			entity.MethodDelete,
			args{
				i: []byte(`"DELETE"`),
			},
			false,
		},
		{
			"MethodConnect UnmarshalJSON()",
			entity.MethodConnect,
			args{
				i: []byte(`"CONNECT"`),
			},
			false,
		},
		{
			"MethodOptions UnmarshallJSON()",
			entity.MethodOptions,
			args{
				i: []byte(`"OPTIONS"`),
			},
			false,
		},
		{
			"MethodTrace UnmarshallJSON()",
			entity.MethodTrace,
			args{
				i: []byte(`"TRACE"`),
			},
			false,
		},
		{
			"MethodGet UnmarshallJSON() prefix not found error",
			entity.MethodGet,
			args{
				i: []byte(`GET"`),
			},
			true,
		},
		{
			"MethodGet UnmarshallJSON() suffix not found error",
			entity.MethodGet,
			args{
				i: []byte(`"GET`),
			},
			true,
		},
		{
			"undefined method UnmarshallJSON()",
			entity.MethodGet,
			args{
				i: []byte(`"undefined"`),
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
