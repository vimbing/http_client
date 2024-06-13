package httpv3

import (
	"context"
	"io"

	fhttp "github.com/vimbing/fhttp"
)

type requestExecuteResult struct {
	response *fhttp.Response
	err      error
}

func (c *Client) executeRequest(req *Request, resultChan chan *requestExecuteResult) {
	defer close(resultChan)

	if err := req.Build(); err != nil {
		resultChan <- &requestExecuteResult{err: err}
		return
	}

	fhttpRes, err := c.fhttpClient.Do(req.fhttpRequest)
	resultChan <- &requestExecuteResult{err: err, response: fhttpRes}
}

func (c *Client) Do(req *Request) (*Response, error) {
	for _, m := range c.cfg.requestMiddleware {
		if err := m(req); err != nil {
			return &Response{}, err
		}
	}

	ctx, cancel := context.WithTimeout(context.Background(), c.cfg.timeout)

	defer cancel()

	resultChan := make(chan *requestExecuteResult, 1)

	go c.executeRequest(req, resultChan)

	var fhttpRes *fhttp.Response

	select {
	case result := <-resultChan:
		if result == nil {
			return &Response{}, ErrResponseNil
		}

		if result.err != nil {
			return &Response{}, result.err
		}

		fhttpRes = result.response
	case <-ctx.Done():
		return &Response{}, ErrRequestTimedOut
	}

	defer fhttpRes.Body.Close()

	body, err := io.ReadAll(fhttpRes.Body)

	if err != nil {
		return &Response{}, err
	}

	decodedBody, err := decodeResponseBody(fhttpRes.Header, body)

	if err != nil {
		return &Response{}, err
	}

	res := &Response{
		body:          decodedBody,
		fhttpResponse: fhttpRes,
	}

	for _, m := range c.cfg.responseMiddleware {
		if err := m(res); err != nil {
			return &Response{}, err
		}
	}

	return res, nil
}

func (c *Client) UseRequest(f RequestMiddlewareFunc) {
	c.cfg.requestMiddleware = append(c.cfg.requestMiddleware, f)
}

func (c *Client) UseResponse(f ResponseMiddlewareFunc) {
	c.cfg.responseMiddleware = append(c.cfg.responseMiddleware, f)
}
