package http_client

import (
	"fmt"
	"net/url"
	"testing"

	"github.com/vimbing/fhttp/cookiejar"
)

func TestCookieGetsSet(t *testing.T) {
	testCases := []struct {
		cookieName  string
		cookieValue string
	}{
		{cookieName: "foo", cookieValue: "bar"},
		{cookieName: "bar", cookieValue: "foo"},
		{cookieName: "test_cookie", cookieValue: "test_value"},
		{cookieName: "cookie1", cookieValue: "cookie1_value"},
	}

	expectedCookiesCount := 1

	for _, testCase := range testCases {
		jar, _ := cookiejar.New(nil)

		client := MustNew(
			WithCookieJar(jar),
		)

		urlString := fmt.Sprintf(
			"http://127.0.0.1:%d/cookie-set?cookieName=%s&cookieValue=%s",
			testServerPort,
			testCase.cookieName,
			testCase.cookieValue,
		)

		res, err := client.Get(urlString)

		if err != nil {
			t.Fatalf("Unexpected error while pinging test server: %v", err)
		}

		if res.StatusCode() != 200 {
			t.Fatalf("Unexpected status code from test server: %s", res.Status())
		}

		parsedUrl, err := url.Parse(urlString)

		if err != nil {
			t.Fatalf("Unexpected error while parsing url of test server: %v", err)
		}

		cookies := client.GetCookies(parsedUrl)

		if len(cookies) != expectedCookiesCount {
			t.Fatalf("Unexpected length of cookies in jar got %d, expected: %d", len(cookies), expectedCookiesCount)
		}

		cookie := cookies[0]

		if cookie.Name != testCase.cookieName {
			t.Fatalf("Unexpected cookie name got %s, expected: %s", cookie.Name, testCase.cookieName)
		}

		if cookie.Value != testCase.cookieValue {
			t.Fatalf("Unexpected cookie name got %s, expected: %s", cookie.Value, testCase.cookieValue)
		}
	}
}

func TestCookieAdd(t *testing.T) {
	testCases := []struct {
		cookieName  string
		cookieValue string
	}{
		{cookieName: "foo", cookieValue: "bar"},
		{cookieName: "bar", cookieValue: "foo"},
		{cookieName: "test_cookie", cookieValue: "test_value"},
		{cookieName: "cookie1", cookieValue: "cookie1_value"},
	}

	expectedCookiesCount := 1

	for _, testCase := range testCases {
		jar, _ := cookiejar.New(nil)

		client := MustNew(
			WithCookieJar(jar),
		)

		parsedUrl, err := url.Parse(fmt.Sprintf("http://127.0.0.1:%d/cookie-set", testServerPort))

		if err != nil {
			t.Fatalf("Unexpected error while parsing test url: %v", err)
		}

		client.AddCookieSimple(parsedUrl, testCase.cookieName, testCase.cookieValue)

		cookies := client.GetCookies(parsedUrl)

		if len(cookies) != expectedCookiesCount {
			t.Fatalf("Unexpected length of cookies in jar got %d, expected: %d", len(cookies), expectedCookiesCount)
		}

		cookie := cookies[0]

		if cookie.Name != testCase.cookieName {
			t.Fatalf("Unexpected cookie name got %s, expected: %s", cookie.Name, testCase.cookieName)
		}

		if cookie.Value != testCase.cookieValue {
			t.Fatalf("Unexpected cookie name got %s, expected: %s", cookie.Value, testCase.cookieValue)
		}
	}
}

func TestCookieDelete(t *testing.T) {
	testCases := []struct {
		cookieName  string
		cookieValue string
	}{
		{cookieName: "foo", cookieValue: "bar"},
		{cookieName: "bar", cookieValue: "foo"},
		{cookieName: "test_cookie", cookieValue: "test_value"},
		{cookieName: "cookie1", cookieValue: "cookie1_value"},
	}

	expectedCookiesCount := 0

	for _, testCase := range testCases {
		jar, _ := cookiejar.New(nil)

		client := MustNew(
			WithCookieJar(jar),
		)

		parsedUrl, err := url.Parse(fmt.Sprintf("http://127.0.0.1:%d/cookie-set", testServerPort))

		if err != nil {
			t.Fatalf("Unexpected error while parsing test url: %v", err)
		}

		client.AddCookieSimple(parsedUrl, testCase.cookieName, testCase.cookieValue)
		client.RemoveCookie(parsedUrl, testCase.cookieName)

		cookies := client.GetCookies(parsedUrl)

		if len(cookies) != expectedCookiesCount {
			t.Fatalf("Unexpected length of cookies in jar got %d, expected: %d", len(cookies), expectedCookiesCount)
		}
	}
}

func TestCookieUpdate(t *testing.T) {
	testCases := []struct {
		cookieName   string
		initialValue string
		updateValue  string
	}{
		{cookieName: "foo", initialValue: "bar", updateValue: "updated"},
		{cookieName: "bar", initialValue: "foo", updateValue: "updated"},
		{cookieName: "test_cookie", initialValue: "test_value", updateValue: "updated"},
		{cookieName: "cookie1", initialValue: "cookie1_value", updateValue: "updated"},
	}

	expectedCookiesCount := 1

	for _, testCase := range testCases {
		jar, _ := cookiejar.New(nil)

		client := MustNew(
			WithCookieJar(jar),
		)

		parsedUrl, err := url.Parse(fmt.Sprintf("http://127.0.0.1:%d/cookie-set", testServerPort))

		if err != nil {
			t.Fatalf("Unexpected error while parsing test url: %v", err)
		}

		client.AddCookieSimple(parsedUrl, testCase.cookieName, testCase.initialValue)
		client.UpdateCookieSimple(parsedUrl, testCase.cookieName, testCase.updateValue)

		cookies := client.GetCookies(parsedUrl)

		if len(cookies) != expectedCookiesCount {
			t.Fatalf("Unexpected length of cookies in jar got %d, expected: %d", len(cookies), expectedCookiesCount)
		}

		cookie := cookies[0]

		if cookie.Value != testCase.updateValue {
			t.Fatalf("Unexpected cookie value after update got %s, expected: %s", cookie.Value, testCase.updateValue)
		}
	}
}
