package errors

import (
	"fmt"
	"net/http"
)

const (
	// IMPORTANT NOTE PLEASE DON'T CHANGE THE ORDER OF THE ENUMS
	errCodeInvalidArgument = iota
	errCodeRecordNotFound
	errCodeAlreadyExists
	errCodeInternalServer
	errCodeBadRequest
	errCodeUnknown
	errCodePermissionDenied
	errFailedPrecondition
	errResourceExhausted
)

// ErrorResponse represents a generic error response structure
type ErrorResponse struct {
	Message      string            `json:"message"`
	DebugMessage string            `json:"debug_message,omitempty"`
	Code         int               `json:"code"`
	ErrorType    ErrorType         `json:"error_type"`
	ErrorDetails map[string]string `json:"error_details,omitempty"`
}

func (e ErrorResponse) Error() string {
	if e.DebugMessage == "" {
		return fmt.Sprintf("message: %s, code: %d", e.Message, e.Code)
	}

	return fmt.Sprintf("message: %s, debug_message: %s, code: %d", e.Message, e.DebugMessage, e.Code)
}

func NewErrorResponse(err error, errType ErrorType) ErrorResponse {
	return ErrorResponse{
		Code:      GetErrorCodeForErrorType(errType),
		Message:   err.Error(),
		ErrorType: errType,
	}
}

func NewErrorResponseWithCode(msg string, debugMsg string, code int) ErrorResponse {
	return ErrorResponse{
		Message:      msg,
		Code:         code,
		DebugMessage: debugMsg,
	}
}

func NewErrorResponseWithDebug(msg string, debug string, errType ErrorType) ErrorResponse {
	return ErrorResponse{
		Code:         GetErrorCodeForErrorType(errType),
		Message:      msg,
		DebugMessage: debug,
		ErrorType:    errType,
	}
}

func GetHttpCodeFromErrorType(errType ErrorType) int {
	switch errType {
	case ErrInvalidArgumentStr:
		return http.StatusBadRequest
	case ErrRecordNotFoundStr:
		return http.StatusOK
	case ErrAlreadyExistsStr:
		return http.StatusAlreadyReported
	case ErrInternalServerStr:
		return http.StatusInternalServerError
	case ErrBadRequestStr:
		return http.StatusBadRequest
	case ErrPermissionDeniedStr:
		return http.StatusForbidden
	case ErrFailedPreconditionStr:
		return http.StatusPreconditionFailed
	case ErrResourceExhaustedStr:
		return http.StatusTooManyRequests
	default:
		return http.StatusInternalServerError
	}
}

func GetErrorCodeForErrorType(errType ErrorType) int {
	switch errType {
	case ErrInvalidArgumentStr:
		return errCodeInvalidArgument
	case ErrRecordNotFoundStr:
		return errCodeRecordNotFound
	case ErrAlreadyExistsStr:
		return errCodeAlreadyExists
	case ErrInternalServerStr:
		return errCodeInternalServer
	case ErrBadRequestStr:
		return errCodeBadRequest
	case ErrPermissionDeniedStr:
		return errCodePermissionDenied
	case ErrFailedPreconditionStr:
		return errFailedPrecondition
	case ErrResourceExhaustedStr:
		return errResourceExhausted
	default:
		return errCodeUnknown
	}
}

func GetErrorTypeFromErrorCode(code int) ErrorType {
	switch code {
	case http.StatusBadRequest, errCodeBadRequest:
		return ErrBadRequestStr
	case errCodeInvalidArgument:
		return ErrInvalidArgumentStr
	case http.StatusNotFound, errCodeRecordNotFound:
		return ErrRecordNotFoundStr
	case http.StatusAlreadyReported, errCodeAlreadyExists:
		return ErrAlreadyExistsStr
	case http.StatusInternalServerError, errCodeInternalServer:
		return ErrInternalServerStr
	case errCodePermissionDenied:
		return ErrPermissionDeniedStr
	case errFailedPrecondition:
		return ErrFailedPreconditionStr
	case errResourceExhausted:
		return ErrResourceExhaustedStr
	default:
		return ErrInternalServerStr
	}
}
