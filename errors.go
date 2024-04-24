package main

import (
	"database/sql"
	"errors"
	"github.com/go-sql-driver/mysql"
	"net/http"
)

type ApiError struct {
	statusCode int
	message    string
}

func (e *ApiError) Error() string {
	return e.message
}

func (e *ApiError) ApiError() (int, string) {
	return e.statusCode, e.message
}

func WrapError(err error) *ApiError {
	switch {
	case checkIfDuplicateError(err):
		return &ApiError{statusCode: http.StatusConflict, message: err.Error()}
	case errors.Is(err, sql.ErrNoRows):
		// handle no DB rows found
		return &ApiError{statusCode: http.StatusNotFound, message: err.Error()}
	case errors.Is(err, NotEnoughBalanceErr):
		return &ApiError{statusCode: http.StatusBadRequest, message: err.Error()}
	default:
		// return the standard
		return &ApiError{statusCode: http.StatusInternalServerError, message: "Internal server error"}
	}
}

func checkIfDuplicateError(err error) bool {
	var mysqlErr *mysql.MySQLError
	return errors.As(err, &mysqlErr) && mysqlErr.Number == 1062
}

func handleApiError(w http.ResponseWriter, err error) {
	if err, ok := err.(*ApiError); ok {
		statusCode, msg := err.ApiError()
		http.Error(w, msg, statusCode)
	} else {
		http.Error(w, "Unknown error", http.StatusInternalServerError)
	}
}
