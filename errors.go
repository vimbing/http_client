package httpv3

import "errors"

var (
	ErrResponseNil     = errors.New("request ended up with nil response")
	ErrRequestTimedOut = errors.New("request timed out")
)
