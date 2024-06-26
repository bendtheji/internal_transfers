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

var (
	NotEnoughBalanceErr         = errors.New("not enough balance")
	InvalidAccountIDErr         = errors.New("invalid account id")
	InvalidInitialBalanceErr    = errors.New("invalid initial balance")
	InvalidTransactionAmountErr = errors.New("invalid transaction amount")
	ReqUnmarshalTypeErr         = errors.New("invalid request type")
)

func HandleError(err error) *ApiError {
	var mysqlErr *mysql.MySQLError
	switch {
	case errors.As(err, &mysqlErr):
		return &ApiError{statusCode: getStatusCodeForMySQLErr(mysqlErr.Number), message: err.Error()}

	case errors.Is(err, sql.ErrNoRows):
		return &ApiError{statusCode: http.StatusNotFound, message: err.Error()}

	case errors.Is(err, NotEnoughBalanceErr):
		return &ApiError{statusCode: http.StatusBadRequest, message: err.Error()}
	case errors.Is(err, InvalidAccountIDErr):
		return &ApiError{statusCode: http.StatusBadRequest, message: err.Error()}
	case errors.Is(err, InvalidInitialBalanceErr):
		return &ApiError{statusCode: http.StatusBadRequest, message: err.Error()}
	case errors.Is(err, InvalidTransactionAmountErr):
		return &ApiError{statusCode: http.StatusBadRequest, message: err.Error()}
	case errors.Is(err, ReqUnmarshalTypeErr):
		return &ApiError{statusCode: http.StatusBadRequest, message: err.Error()}

	case errors.Is(err, context.DeadlineExceeded):
		return &ApiError{statusCode: http.StatusGatewayTimeout, message: err.Error()}

	default:
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
