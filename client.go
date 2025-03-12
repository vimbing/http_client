package http

import (
	"bytes"
	"errors"
	"io"

	"github.com/samber/lo"
	"github.com/vimbing/fhttp/cookiejar"
)

func (c *Client) BindJar(jar *cookiejar.Jar) {
	c.cfg.jar = jar
	c.fhttpClient.Jar = jar
}

func (c *Client) executeRequest(req *Request, resultChan chan *requestExecutionResult) {
	defer close(resultChan)

	fhttpRes, err := c.fhttpClient.Do(req.fhttpRequest)

	if err != nil {
		if len(c.cfg.responseErrorMiddleware) > 0 {
			for _, m := range c.cfg.responseErrorMiddleware {
				m(req, err)
			}
		}

		resultChan <- &requestExecutionResult{
			error: err,
		}

		return
	}

	defer fhttpRes.Body.Close()

	decodedBody, err := decodeResponseBody(fhttpRes.Header, fhttpRes.Body)

	if err != nil {
		resultChan <- &requestExecutionResult{
			error: err,
		}

		return
	}

	buff := bytes.NewBuffer([]byte{})
	defer buff.Reset()

	if _, err := io.Copy(buff, decodedBody); err != nil {
		resultChan <- &requestExecutionResult{
			error: err,
		}

		return
	}

	res := &Response{
		Body:          buff.Bytes(),
		fhttpResponse: fhttpRes,
	}

	for _, m := range c.cfg.responseMiddleware {
		if err := m(res); err != nil {
			resultChan <- &requestExecutionResult{
				res:   res,
				error: err,
			}

			return
		}
	}

	resultChan <- &requestExecutionResult{
		res:   res,
		error: err,
	}
}

func (c *Client) RotateProxy() error {
	return rebindRoundtripper(c.fhttpClient, c.cfg)
}

func (c *Client) Do(req *Request) (*Response, error) {
	for _, m := range c.cfg.requestMiddleware {
		if err := m(req); err != nil {
			return &Response{}, err
		}
	}

	req.tlsProfile = c.cfg.tlsProfile

	ctx, reqCtxCancel, err := req.Build(
		c.cfg.timeout,
	)

	defer reqCtxCancel()

	if err != nil {
		return &Response{}, err
	}

	resultChan := make(chan *requestExecutionResult, 1)

	if c.cfg.forceRotation {
		err := rebindRoundtripper(c.fhttpClient, c.cfg)

		if err != nil {
			return nil, err
		}
	}

	go c.executeRequest(req, resultChan)

	for {
		select {
		case result := <-resultChan:
			if result.error != nil {
				return result.res, result.error
			}

			if c.cfg.statusValidationFunc != nil {
				return result.res, c.cfg.statusValidationFunc(result.res.StatusCode(), c)
			}

			return result.res, nil
		case <-ctx.Done():
			return nil, errors.New("context cancelled")
		}
	}
}

func (c *Client) UseRequest(f RequestMiddlewareFunc) {
	c.cfg.requestMiddleware = append(c.cfg.requestMiddleware, f)
}

func (c *Client) UseResponse(f ResponseMiddlewareFunc) {
	c.cfg.responseMiddleware = append(c.cfg.responseMiddleware, f)
}

func (c *Client) UseResponseError(f ResponseErrorMiddlewareFunc) {
	c.cfg.responseErrorMiddleware = append(c.cfg.responseErrorMiddleware, f)
}

func (c *Client) ChangeProxy(proxy string) {
	p, _ := parseSingleProxy(proxy)
	c.cfg.proxies = []string{string(p)}
	rebindRoundtripper(c.fhttpClient, c.cfg)
}

func (c *Client) ChangeProxyList(proxies []string) {
	c.cfg.proxies = lo.Map(parseList(proxies), func(p OptionProxy, i int) string { return string(p) })
	rebindRoundtripper(c.fhttpClient, c.cfg)
}

func (c *Client) ChangeProxyParsed(proxy string) {
	c.cfg.proxies = []string{proxy}
	rebindRoundtripper(c.fhttpClient, c.cfg)
}

func (c *Client) ChangeProxyListParsed(proxies []string) {
	c.cfg.proxies = proxies
	rebindRoundtripper(c.fhttpClient, c.cfg)
}
