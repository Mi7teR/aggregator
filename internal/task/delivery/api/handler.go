package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/Mi7teR/aggregator/internal/task/entity"
	"github.com/Mi7teR/aggregator/internal/task/repository"
	"github.com/Mi7teR/aggregator/internal/task/service"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type Handler struct {
	s *service.Service
}

func NewHandler(s *service.Service) *Handler {
	return &Handler{s: s}
}

func (h *Handler) AddTask(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "application/json")

	var req entity.Task
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(&entity.ErrorResponse{Error: err.Error()})
		return
	}

	taskID, err := h.s.AddTask(r.Context(), &req)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(&entity.ErrorResponse{Error: err.Error()})
		return
	}

	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(&entity.TaskResult{ID: taskID})
}

func (h *Handler) GetTaskResult(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "application/json")

	id := chi.URLParam(r, "id")
	if _, err := uuid.Parse(id); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(&entity.ErrorResponse{Error: fmt.Errorf("uuid parse: %w", err).Error()})
		return
	}

	res, err := h.s.GetTaskResult(r.Context(), id)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			w.WriteHeader(http.StatusNotFound)
			_ = json.NewEncoder(w).Encode(&entity.ErrorResponse{Error: err.Error()})
			return
		}

		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(&entity.ErrorResponse{Error: err.Error()})
		return
	}

	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(res)
}

func (h *Handler) NotFoundHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "application/json")
	w.WriteHeader(http.StatusNotFound)
	_ = json.NewEncoder(w).Encode(&entity.ErrorResponse{Error: "page not found"})
}

func (h *Handler) MethodNotAllowedHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "application/json")
	w.WriteHeader(http.StatusMethodNotAllowed)
	_ = json.NewEncoder(w).Encode(&entity.ErrorResponse{Error: "method not allowed"})
}
