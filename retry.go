package http

import (
	"slices"
	"time"
)

func (r *Retry) Retry(f doFunc, req *Request) (*Response, error) {
	if r.Max == 0 {
		return f(req)
	}

	for i := 0; i < r.Max; i++ {
		if i != 0 {
			time.Sleep(r.Delay)
		}

		res, err := f(req)

		if err == nil {
			return res, nil
		}

		if slices.Contains(r.EndingErrors, err) {
			return res, err
		}

		if slices.Contains(r.IgnoredErrors, err) {
			i--
		}
	}

	return nil, ErrRetryExceed
}
