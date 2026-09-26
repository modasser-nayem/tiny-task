package apperrors

import "net/http"

type AppError struct {
	Code    string
	Message string
	Status  int
}

func (e *AppError) Error() string {
	return e.Message
}

// Bad  Request
func BadRequest(code, message string) *AppError {
	return &AppError{
		Code:    code,
		Message: message,
		Status:  http.StatusBadRequest,
	}
}

// Unauthorized
func Unauthorized(code, message string) *AppError {
	return &AppError{
		Code:    code,
		Message: message,
		Status:  http.StatusUnauthorized,
	}
}

// Forbidden
func Forbidden(code, message string) *AppError {
	return &AppError{
		Code:    code,
		Message: message,
		Status:  http.StatusForbidden,
	}
}

// Not Found
func NotFound(code, message string) *AppError {
	return &AppError{
		Code:    code,
		Message: message,
		Status:  http.StatusNotFound,
	}
}

// Conflict
func Conflict(code, message string) *AppError {
	return &AppError{
		Code:    code,
		Message: message,
		Status:  http.StatusConflict,
	}
}


// Internal Server Error
func Internal(code, message string) *AppError {
	return &AppError{
		Code:    code,
		Message: message,
		Status:  http.StatusInternalServerError,
	}
}

