package utils

import (
	"database/sql"
	"errors"
	"net/http"
)

const (
	ValidationErrorMessage string = "validation error: incorrect params"
	NotFoundMessage        string = "not found"
	InternalErrorMessage   string = "internal error"
)

type ErrorResult struct {
	Err        error
	Msg        string
	StatusCode int
}

func (e *ErrorResult) Error() string {
	return e.Err.Error()
}

func WrapBadRequestError(err error) *ErrorResult {
	return &ErrorResult{
		Err:        err,
		Msg:        ValidationErrorMessage,
		StatusCode: http.StatusBadRequest,
	}
}

func WrapForbiddenError(err error, msg string) *ErrorResult {
	return &ErrorResult{
		Err:        err,
		Msg:        msg,
		StatusCode: http.StatusForbidden,
	}
}

func WrapNotFoundError(err error, msg string) *ErrorResult {
	return &ErrorResult{
		Err:        err,
		Msg:        msg,
		StatusCode: http.StatusNotFound,
	}
}

func WrapInternalError(err error) *ErrorResult {
	return &ErrorResult{
		Err:        err,
		Msg:        InternalErrorMessage,
		StatusCode: http.StatusInternalServerError,
	}
}

func WrapError(err error, msg string, status int) *ErrorResult {
	return &ErrorResult{
		Err:        err,
		Msg:        msg,
		StatusCode: status,
	}
}

func FromError(err error) (*ErrorResult, bool) {
	if err == nil {
		return nil, false
	}

	var result *ErrorResult
	ok := errors.As(err, &result)
	if !ok {
		return nil, false
	}

	return result, true
}

func WrapSqlError(err error) error {
	switch {
	case errors.Is(err, sql.ErrNoRows):
		return WrapNotFoundError(err, NotFoundMessage)
	default:
		return WrapInternalError(err)
	}
}
