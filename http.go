package http

import (
	"fmt"
	"io"
	urlLib "net/url"
	"strings"

	http "github.com/vimbing/fhttp"
	"github.com/vimbing/retry"
)

func (c *Client) reinitFhttpClient() error {
	if c.cfg.tlsProfile != nil && c.cfg.tlsProfile.Ja3.IsSet() {
		c.cfg.ja3 = c.cfg.tlsProfile.Ja3
	}

	newClient, err := newFhttpClient(c.cfg)

	if err != nil {
		return err
	}

	c.fhttpClient = newClient

	if c.fhttpClient != nil && c.cfg.jar != nil {
		c.fhttpClient.Jar = c.cfg.jar
	}

	return err
}

func (c *Client) urlValues(v map[string]string) urlLib.Values {
	values := urlLib.Values{}

	for key, value := range v {
		values.Add(key, value)
	}

	return values
}

func (c *Client) NewRequest(url string, options ...any) (*Request, error) {
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
		case QueryParams:
			values := c.urlValues(v)

			var joinChar = "?"

			if strings.Contains(req.Url, "?") {
				joinChar = "&"
			}

			req.Url = fmt.Sprintf("%s%s%s", req.Url, joinChar, values.Encode())
		case FormUrlEncoded:
			req.Body = strings.NewReader(c.urlValues(v).Encode())

			if len(req.Header.Get("content-type")) == 0 {
				req.Header.Set("content-type", "application/x-www-form-urlencoded")
			}
		case retry.Retrier:
			req.retrier = &v
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
	req, err := c.NewRequest(url, options...)

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
