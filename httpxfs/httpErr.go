package httpxfs

import "fmt"

type RPCError interface {
	error
	Code() int
	Message() string
}
type rpcError struct {
	Code    int    `json:"code,omitempty"`
	Message string `json:"message,omitempty"`
}

var (
	parseError          = NewRPCError(-32700, "Parse error")
	invalidRequestError = NewRPCError(-32600, "Invalid request")
	methodNotFoundError = NewRPCError(-32601, "Method not found")
	invalidParamsError  = NewRPCError(-32602, "Invalid params")
	internalError       = NewRPCError(-32603, "Internal error")
	RequireParamError   = func(msg string, params ...interface{}) *rpcError {
		return NewRPCError(-32000, fmt.Sprintf(msg, params...))
	}
	ParamsParseError = func(msg string, params ...interface{}) *rpcError {
		return NewRPCError(-32001, fmt.Sprintf(msg, params...))
	}
	LoadStateTreeError = func(msg string, params ...interface{}) *rpcError {
		return NewRPCError(-32002, fmt.Sprintf(msg, params...))
	}
)

func NewRPCError(code int, message string) *rpcError {
	return &rpcError{
		Code:    code,
		Message: message,
	}
}

func NewRPCErrorCause(code int, err error) *rpcError {
	return &rpcError{
		Code:    code,
		Message: err.Error(),
	}
}

func (e *rpcError) Error() string {
	return fmt.Sprintf("%d: %s", e.Code, e.Message)
}
