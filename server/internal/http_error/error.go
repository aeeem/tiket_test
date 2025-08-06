package http_error

import (
	"errors"
	"fmt"

	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
)

type HttpErrorCode struct {
	HTTPErrorCode int
	Message       string
}

func CheckError(err error) (httpErr HttpErrorCode) {

	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		switch pgErr.Code {
		case pgerrcode.UniqueViolation:
			httpErr.HTTPErrorCode = 400
			httpErr.Message = fmt.Sprintf("Item with name %s already exists", pgErr.ConstraintName)
			return
		}
	} else {
		httpErr.HTTPErrorCode = 500
		httpErr.Message = err.Error()
		return
	}
	return
}
