package httpv3

import "errors"

var (
	ErrResponseNil          = errors.New("request ended up with nil response")
	ErrRequestTimedOut      = errors.New("request timed out")
	ErrProxyFormatCorrupted = errors.New("proxy format corrupted, cannot parse")
)
