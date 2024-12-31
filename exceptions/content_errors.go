package exceptions

import "errors"

var (
	ErrContentNotFound = errors.New("entities not found")
)
