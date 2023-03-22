package entity

import "net/http"

type TaskResult struct {
	ID             string           `json:"id"`
	Status         TaskResultStatus `json:"status,omitempty"`
	HTTPStatusCode int              `json:"httpStatusCode,omitempty"`
	Headers        http.Header      `json:"headers,omitempty"`
	Length         int64            `json:"length,omitempty"`
}
