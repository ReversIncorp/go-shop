package utils

import (
	"errors"
	errorhandling "marketplace/pkg/error_handling"

	"github.com/ztrue/tracerr"
)

func IsHttpError(err error) bool {
	var appErr *errorhandling.ResponseError
	return errors.As(err, &appErr)
}

func GetHttpErrorOrTracerrError(err error) error {
	if IsHttpError(err) {
		return err
	} else {
		return tracerr.Wrap(err)
	}
}
