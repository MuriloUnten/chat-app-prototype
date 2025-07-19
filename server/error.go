package main

import (
	"fmt"
	"net/http"
)

type APIError struct {
	StatusCode int
	Msg string
}

func (e APIError) Error() string {
	return fmt.Sprintf("api error: %d", e.StatusCode)
}

func NewAPIError(status int, msg string) APIError {
	return APIError{
		StatusCode: status,
		Msg: msg,
	}
}

func InternalError() APIError {
	return APIError{
		StatusCode: http.StatusInternalServerError,
		Msg: "internal error",
	}
}

func BadRequest() APIError {
	return APIError{
		StatusCode: http.StatusBadRequest,
		Msg: "bad request",
	}
}
