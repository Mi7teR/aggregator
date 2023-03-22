package entity

import (
	"bytes"
	"errors"
)

type TaskResultStatus int

const (
	TaskStatusNew TaskResultStatus = iota + 1
	TaskStatusInProcess
	TaskStatusError
	TaskStatusDone
)

var ErrInvalidStatus = errors.New("invalid status")

func (t *TaskResultStatus) UnmarshalJSON(i []byte) error {
	var status TaskResultStatus

	i, ok := bytes.CutPrefix(i, []byte("\""))
	if !ok {
		return ErrPrefixNotFound
	}

	i, ok = bytes.CutSuffix(i, []byte("\""))
	if !ok {
		return ErrSuffixNotFound
	}

	switch string(i) {
	case "new":
		status = TaskStatusNew
	case "in_process":
		status = TaskStatusInProcess
	case "error":
		status = TaskStatusError
	case "done":
		status = TaskStatusDone
	default:
		return ErrInvalidStatus
	}

	*t = status

	return nil
}

func (t *TaskResultStatus) MarshalJSON() ([]byte, error) {
	if *t > TaskStatusDone || *t < TaskStatusNew {
		return nil, ErrInvalidStatus
	}

	b := bytes.Buffer{}

	b.WriteByte('"')
	b.WriteString(t.String())
	b.WriteByte('"')

	return b.Bytes(), nil
}

func (t *TaskResultStatus) String() string {
	var status string

	switch *t {
	case TaskStatusNew:
		status = "new"
	case TaskStatusInProcess:
		status = "in_process"
	case TaskStatusError:
		status = "error"
	case TaskStatusDone:
		status = "done"
	}

	return status
}
