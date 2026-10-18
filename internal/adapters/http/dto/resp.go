package dto

import "github.com/loopopen/pkg/errorx"

type Resp struct {
	Data     any               `json:"data,omitempty"`
	Reason   string            `json:"reason,omitempty"`
	Message  string            `json:"message,omitempty"`
	Error    string            `json:"error,omitempty"`
	Metadata map[string]string `json:"metadata,omitempty"`
}

type Response[T any] struct {
	Data T `json:"data"`
}

type ErrorResponse struct {
	Reason   string            `json:"reason,omitempty"`
	Message  string            `json:"message,omitempty"`
	Cause    string            `json:"error,omitempty"`
	Metadata map[string]string `json:"metadata,omitempty"`
	status   int
}

func (e *ErrorResponse) Error() string {
	if e.Cause != "" {
		return e.Cause
	}
	return e.Message
}

func (e *ErrorResponse) GetStatus() int {
	return e.status
}

func NewError(status int, reason, message, cause string, metadata map[string]string) *ErrorResponse {
	return &ErrorResponse{
		Reason:   reason,
		Message:  message,
		Cause:    cause,
		Metadata: metadata,
		status:   status,
	}
}

func Err(err *errorx.Error) *ErrorResponse {
	var cause string
	if unwrapped := err.Unwrap(); unwrapped != nil {
		cause = unwrapped.Error()
	}
	return NewError(int(err.Code), err.Reason, err.Message, cause, err.Metadata)
}

func OK[T any](data T) Response[T] {
	return Response[T]{Data: data}
}
