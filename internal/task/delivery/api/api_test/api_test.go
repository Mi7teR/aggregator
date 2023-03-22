package api_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"sync"
	"testing"
	"time"

	"github.com/Mi7teR/aggregator/internal/task/delivery/api"
	"github.com/Mi7teR/aggregator/internal/task/entity"
	"github.com/Mi7teR/aggregator/internal/task/repository"
	"github.com/Mi7teR/aggregator/internal/task/service"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

func executeRequest(r *http.Request, router *chi.Mux) *httptest.ResponseRecorder {
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, r)
	return rr
}

func checkResponseCode(t *testing.T, expected, actual int) {
	if expected != actual {
		t.Errorf("Expected response code %d. Got %d\n", expected, actual)
	}
}

func TestTask(t *testing.T) {
	repo := repository.NewTaskInMemoryRepository()
	timeout := time.Second * 1000
	s := service.NewService(repo, timeout)
	h := api.NewHandler(s)
	r := api.NewRouter(h)

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

	var id string
	t.Run("add task", func(t *testing.T) {
		task := &entity.Task{
			Method: entity.MethodPost,
			URL:    fmt.Sprintf("%s/test-path", server.URL),
			Headers: map[string]string{
				"Accept": "application/json",
			},
		}

		taskBody, err := json.Marshal(&task)
		if err != nil {
			t.Errorf("expected to marshal json, got %v", err)
		}

		req, err := http.NewRequest(http.MethodPost, "/task", bytes.NewReader(taskBody))
		if err != nil {
			t.Errorf("expected to create request, got %v", err)
		}

		res := executeRequest(req, r)

		checkResponseCode(t, http.StatusOK, res.Code)

		var responseJSON entity.TaskResult

		err = json.NewDecoder(res.Body).Decode(&responseJSON)
		if err != nil {
			t.Errorf("expected to decode response, got %v", err)
		}

		_, err = uuid.Parse(responseJSON.ID)
		if err != nil {
			t.Errorf("expected to parse valid uuid, got %v", err)
		}

		id = responseJSON.ID
	})

	t.Run("get task", func(t *testing.T) {
		wg.Wait()
		req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("/task/%s", id), nil)
		if err != nil {
			t.Errorf("expected to create request, got %v", err)
		}

		res := executeRequest(req, r)

		checkResponseCode(t, http.StatusOK, res.Code)

		var responseJSON entity.TaskResult

		expectedResult := entity.TaskResult{
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

		err = json.NewDecoder(res.Body).Decode(&responseJSON)
		if err != nil {
			t.Errorf("expected to decode response, got %v", err)
		}

		if !reflect.DeepEqual(responseJSON, expectedResult) {
			t.Errorf("GetTask() error, got %v, want %v", responseJSON, expectedResult)
		}
	})
}
