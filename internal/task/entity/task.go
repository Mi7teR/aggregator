package entity

type Task struct {
	Method  TaskMethod        `json:"method"`
	URL     string            `json:"url"`
	Headers map[string]string `json:"headers"`
}
