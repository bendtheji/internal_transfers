package errors

import "net/http"

const (
	MySQLDuplicateEntryErr = 1062
	MySQLDeadLockErr       = 1213
)

func getStatusCodeForMySQLErr(code uint16) int {
	switch code {
	case MySQLDuplicateEntryErr:
		return http.StatusConflict
	case MySQLDeadLockErr:
		return http.StatusServiceUnavailable
	default:
		return http.StatusInternalServerError
	}
}
