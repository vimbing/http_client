package http

import (
	"io"

	"github.com/vimbing/retry"
)

func (c *Client) Do(req *Request) (*Response, error) {
	for _, m := range c.cfg.requestMiddleware {
		if err := m(req); err != nil {
			return &Response{}, err
		}
	}

	req.tlsProfile = c.cfg.tlsProfile

	reqCtxCancel, err := req.Build(
		c.cfg.timeout,
	)

	if err != nil {
		return &Response{}, err
	}

	defer reqCtxCancel()

	if req.retrier == nil {
		req.retrier = &retry.Retrier{Max: 1}
	}

	res := &Response{}

	req.retrier.Retry(func() error {
		fhttpRes, err := c.fhttpClient.Do(req.fhttpRequest)

		if err != nil {
			return err
		}

		defer fhttpRes.Body.Close()

		body, err := io.ReadAll(fhttpRes.Body)

		if err != nil {
			return err
		}

		decodedBody, err := decodeResponseBody(fhttpRes.Header, body)

		if err != nil {
			return err
		}

		res = &Response{
			Body:          decodedBody,
			fhttpResponse: fhttpRes,
		}

		for _, m := range c.cfg.responseMiddleware {
			if err := m(res); err != nil {
				return err
			}
		}

		return nil
	})

	return res, nil
}

func (c *Client) UseRequest(f RequestMiddlewareFunc) {
	c.cfg.requestMiddleware = append(c.cfg.requestMiddleware, f)
}

func (c *Client) UseResponse(f ResponseMiddlewareFunc) {
	c.cfg.responseMiddleware = append(c.cfg.responseMiddleware, f)
}
