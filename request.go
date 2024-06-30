package http

import (
	"context"
	"time"

	fhttp "github.com/vimbing/fhttp"
)

func (r *Request) Build(timeout time.Duration) (context.CancelFunc, error) {
	ctx, cancel := context.WithTimeout(
		context.Background(),
		timeout,
	)

	req, err := fhttp.NewRequestWithContext(ctx, r.Method, r.Url, r.Body)

	if err != nil {
		return cancel, err
	}

	req.Header = r.Header
	r.fhttpRequest = req

	return cancel, nil
}
