package entity

import (
	"bytes"
	"errors"
	"net/http"
)

type TaskMethod int

const (
	MethodGet TaskMethod = iota
	MethodHead
	MethodPost
	MethodPut
	MethodPatch
	MethodDelete
	MethodConnect
	MethodOptions
	MethodTrace
)

var (
	ErrPrefixNotFound = errors.New("prefix not found")
	ErrSuffixNotFound = errors.New("suffix not found")
	ErrInvalidMethod  = errors.New("invalid method")
)

func (t *TaskMethod) UnmarshalJSON(i []byte) error {
	var method TaskMethod

	i, ok := bytes.CutPrefix(i, []byte("\""))
	if !ok {
		return ErrPrefixNotFound
	}

	i, ok = bytes.CutSuffix(i, []byte("\""))
	if !ok {
		return ErrSuffixNotFound
	}

	switch string(i) {
	case http.MethodGet:
		method = MethodGet
	case http.MethodHead:
		method = MethodHead
	case http.MethodPost:
		method = MethodPost
	case http.MethodPut:
		method = MethodPut
	case http.MethodPatch:
		method = MethodPatch
	case http.MethodDelete:
		method = MethodDelete
	case http.MethodConnect:
		method = MethodConnect
	case http.MethodOptions:
		method = MethodOptions
	case http.MethodTrace:
		method = MethodTrace
	default:
		return ErrInvalidMethod
	}

	*t = method

	return nil
}

func (t *TaskMethod) MarshalJSON() ([]byte, error) {
	if *t > MethodTrace || *t < MethodGet {
		return nil, ErrInvalidMethod
	}

	buf := bytes.Buffer{}

	buf.WriteByte('"')
	buf.WriteString(t.String())
	buf.WriteByte('"')

	return buf.Bytes(), nil
}

func (t *TaskMethod) String() string {
	var method string

	switch *t {
	case MethodGet:
		method = http.MethodGet
	case MethodHead:
		method = http.MethodHead
	case MethodPost:
		method = http.MethodPost
	case MethodPut:
		method = http.MethodPut
	case MethodPatch:
		method = http.MethodPatch
	case MethodDelete:
		method = http.MethodDelete
	case MethodConnect:
		method = http.MethodConnect
	case MethodOptions:
		method = http.MethodOptions
	case MethodTrace:
		method = http.MethodTrace
	}

	return method
}
