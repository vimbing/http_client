package http_client

import (
	"net/url"

	http "github.com/vimbing/fhttp"
)

func (c *Client) GetCookies(u *url.URL) []*http.Cookie {
	if c.fhttpClient.Jar == nil {
		return []*http.Cookie{}
	}

	return c.fhttpClient.Jar.Cookies(u)
}

func (c *Client) GetCookieByName(u *url.URL, name string) *http.Cookie {
	if c.fhttpClient.Jar == nil {
		return nil
	}

	cookies := c.GetCookies(u)

	var foundCookie *http.Cookie

	for _, cookie := range cookies {
		if cookie.Name == name {
			foundCookie = cookie
		}
	}

	return foundCookie
}

func (c *Client) RemoveCookie(u *url.URL, name string) {
	cookies := c.GetCookies(u)

	if len(cookies) == 0 {
		return
	}

	for _, cookie := range cookies {
		if cookie.Name == name {
			cookie.MaxAge = -1
		}
	}

	c.fhttpClient.Jar.SetCookies(u, cookies)
}

func (c *Client) AddCookieSimple(u *url.URL, name, value string) {
	c.AddCookie(u, &http.Cookie{Name: name, Value: value})
}

func (c *Client) AddCookie(u *url.URL, cookie *http.Cookie) {
	cookies := []*http.Cookie{}
	cookies = append(cookies, cookie)
	c.fhttpClient.Jar.SetCookies(u, cookies)
}

func (c *Client) UpdateCookieSimple(u *url.URL, name, value string) {
	if c.fhttpClient.Jar == nil {
		return
	}

	c.RemoveCookie(u, name)
	c.AddCookieSimple(u, name, value)
}

func (c *Client) UpdateCookie(u *url.URL, cookie *http.Cookie) {
	if c.fhttpClient.Jar == nil {
		return
	}

	c.RemoveCookie(u, cookie.Name)
	c.AddCookie(u, cookie)
}

func (c *Client) ClearAllCookies(u *url.URL) {
	if c.fhttpClient.Jar == nil {
		return
	}

	cookies := c.GetCookies(u)

	for _, cookie := range cookies {
		c.RemoveCookie(u, cookie.Name)
	}
}
