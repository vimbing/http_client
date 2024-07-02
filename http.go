package http

import (
	"io"

	http "github.com/vimbing/fhttp"
)

func (c *Client) reinitFhttpClient() error {
	var jarCopy http.CookieJar

	if c.fhttpClient != nil {
		jarCopy = c.fhttpClient.Jar
	}

	newClient, err := newFhttpClient(c.cfg)

	if err != nil {
		return err
	}

	if jarCopy != nil {
		newClient.Jar = jarCopy
	}

	c.fhttpClient = newClient

	return err
}

func (c *Client) newRequest(url string, options ...any) (*Request, error) {
	req := &Request{
		Method: "GET",
		Url:    url,
		Body:   nil,
		Header: http.Header{},
	}

	for _, opt := range options {
		switch v := opt.(type) {
		case string:
			req.Method = v
		case http.Header:
			req.Header = v
		case io.Reader:
			req.Body = v
		case RequestJsonBody:
			body, err := marshalAndEncodeBody(v)

			if err != nil {
				return req, err
			}

			req.Body = body
		}
	}

	return req, nil
}

func (c *Client) do(url string, options ...any) (*Response, error) {
	req, err := c.newRequest(url, options...)

	if err != nil {
		return &Response{}, err
	}

	return c.Do(req)
}

func (c *Client) Get(url string, options ...any) (*Response, error) {
	options = append(options, "GET")
	return c.do(url, options...)
}

func (c *Client) Post(url string, options ...any) (*Response, error) {
	options = append(options, "POST")
	return c.do(url, options...)
}

func (c *Client) Put(url string, options ...any) (*Response, error) {
	options = append(options, "PUT")
	return c.do(url, options...)
}

func (c *Client) Delete(url string, options ...any) (*Response, error) {
	options = append(options, "DELETE")
	return c.do(url, options...)
}
