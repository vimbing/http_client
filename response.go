package httpv3

import (
	"fmt"
)

func (r *Response) StatusCode() int {
	return r.fhttpResponse.StatusCode
}

func (r *Response) Status() string {
	return r.fhttpResponse.Status
}

func (r *Response) BodyString() string {
	return string(r.body)
}

func (r *Response) BodyBytes() []byte {
	return r.body
}

// just for logging purposes
func (r *Response) StatusError() error {
	return fmt.Errorf("server responded with unexpected status code: %d", r.fhttpResponse.StatusCode)
}
