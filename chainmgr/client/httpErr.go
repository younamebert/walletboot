package client

import "fmt"

type RPCError struct {
	Code    int    `json:"code,omitempty"`
	Message string `json:"message,omitempty"`
}

func NewRPCError(code int, message string) *RPCError {
	return &RPCError{
		Code:    code,
		Message: message,
	}
}

func NewRPCErrorCause(code int, err error) *RPCError {
	return &RPCError{
		Code:    code,
		Message: err.Error(),
	}
}

func (e *RPCError) Error() string {
	return fmt.Sprintf("%d: %s", e.Code, e.Message)
}
