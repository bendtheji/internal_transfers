package errors

import (
	"context"
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

var NotEnoughBalanceErr = errors.New("Not enough balance")

func WrapError(err error) *ApiError {
	var mysqlErr *mysql.MySQLError
	switch {
	case errors.As(err, &mysqlErr):
		switch mysqlErr.Number {
		case 1062:
			return &ApiError{statusCode: http.StatusConflict, message: err.Error()}
		case 1213:
			return &ApiError{statusCode: http.StatusServiceUnavailable, message: err.Error()}
		default:
			return &ApiError{statusCode: http.StatusInternalServerError, message: err.Error()}
		}
	case errors.Is(err, sql.ErrNoRows):
		return &ApiError{statusCode: http.StatusNotFound, message: err.Error()}
	case errors.Is(err, NotEnoughBalanceErr):
		return &ApiError{statusCode: http.StatusBadRequest, message: err.Error()}
	case errors.Is(err, context.DeadlineExceeded):
		return &ApiError{statusCode: http.StatusGatewayTimeout, message: err.Error()}
	default:
		// return the standard
		return &ApiError{statusCode: http.StatusInternalServerError, message: "Internal server error"}
	}
}

func HandleApiError(w http.ResponseWriter, err error) {
	if err, ok := err.(*ApiError); ok {
		statusCode, msg := err.ApiError()
		http.Error(w, msg, statusCode)
	} else {
		http.Error(w, "Unknown error", http.StatusInternalServerError)
	}
}
