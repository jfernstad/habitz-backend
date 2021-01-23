package endpoints

import (
	"fmt"
	"net/http"
)

// Should be used as parameter for Error() errCode
const (
	BadRequest          = "BAD_REQUEST"
	NotFound            = "NOT_FOUND"
	InternalServerError = "INTERNAL_SERVER_ERROR"
	MissingParameter    = "MISSING_PARAMETER"
	MethodNotAllowed    = "METHOD_NOT_ALLOWED"
	NotImplemented      = "NOT_IMPLEMENTED"
)

// Predefined errors
var (
	ErrEntityNotFound   = newNotFoundErr("request entity not found")
	ErrHandlerNotFound  = newNotFoundErr("handler not found")
	ErrMethodNotAllowed = newError(http.StatusMethodNotAllowed, MethodNotAllowed, "method not allowed")
	ErrNotImplemented   = newError(http.StatusNotImplemented, NotImplemented, "not implemented")
)

// Functions that create API errors by wrapping error object representing the underlying cause.
var (
	WrapJSONDecodeError     = newError(http.StatusBadRequest, BadRequest, "json decode error").Wrap
	WrapInternalServerError = newError(http.StatusInternalServerError, InternalServerError, "internal server error").Wrap
)

func newError(status int, code string, msg string) *errMsg {
	return &errMsg{
		Code:     code,
		HTTPCode: status,
		Message:  msg,
	}
}

// errMsg encodes error json-responses
type errMsg struct {
	Code     string `json:"code"`
	Message  string `json:"message"`
	HTTPCode int    `json:"-"`
}

func (r errMsg) Error() string {
	return fmt.Sprintf("(%d) %s '%s'", r.HTTPCode, r.Code, r.Message)
}

func (r errMsg) String() string {
	return r.Error()
}

func (r *errMsg) Wrap(err error) *errMsg {
	if err == nil { // Make noop if true
		return r
	}

	if r.Message == "" {
		r.Message = err.Error()
	} else {
		r.Message += ": " + err.Error() // This wraps "1st lvl err: 2nd lvl err: 3rd lvl err: ..." Reverse it?
	}
	return r
}

func newBadRequestErr(msg string) *errMsg {
	return &errMsg{
		HTTPCode: http.StatusBadRequest,
		Code:     BadRequest,
		Message:  msg,
	}
}

func newNotFoundErr(msg string) *errMsg {
	return &errMsg{
		HTTPCode: http.StatusNotFound,
		Code:     NotFound,
		Message:  msg,
	}
}

func newInternalServerErr(msg string) *errMsg {
	return &errMsg{
		HTTPCode: http.StatusInternalServerError,
		Code:     InternalServerError,
		Message:  msg,
	}
}

func newMissingParameterErr(msg string) *errMsg {
	return &errMsg{
		HTTPCode: http.StatusBadRequest,
		Code:     MissingParameter,
		Message:  msg,
	}
}
