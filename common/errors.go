package common

import (
	"fmt"
)

type errScope string

const (
	errUnknown errScope = ""
	errConf             = "conf_error"
	errDB               = "db_error"
	errTushare          = "tushare_error"
)

var (
	// Conf error

	// DB error
	ErrDBEmptyUsername = makeError(errDB, "missing Username")
	ErrDBEmptyHost     = makeError(errDB, "missing Host")
	ErrDBEmptyPort     = makeError(errDB, "missing Port")
	ErrDBEmptyDatabase = makeError(errDB, "missing Database")

	// Tushare error
	ErrTushareURL                = makeError(errTushare, "url is empty.")
	ErrTushareToken              = makeError(errTushare, "token is empty.")
	ErrTushareStockFieldsLen     = makeError(errTushare, "length of stock fields is not equal to headers.")
	ErrTushareStockFieldsUnknown = makeError(errTushare, "unknown stock field.")
)

func makeError(scope errScope, msg ...string) error {
	return fmt.Errorf("[%s]: %s", scope, msg)
}
