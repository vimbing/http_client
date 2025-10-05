package http_client

import (
	"errors"
	"fmt"
	"testing"
	"time"

	http "github.com/vimbing/fhttp"
	http2 "github.com/vimbing/fhttp/http2"
	tls "github.com/vimbing/utls"
)

func chrome140Profile() TlsProfile {
	settings := map[http2.SettingID]uint32{
		http2.SettingHeaderTableSize:   65536,
		http2.SettingEnablePush:        0,
		http2.SettingInitialWindowSize: 6291456,
		http2.SettingMaxHeaderListSize: 262144,
	}

	settingsOrder := []http2.SettingID{
		http2.SettingHeaderTableSize,
		http2.SettingEnablePush,
		http2.SettingInitialWindowSize,
		http2.SettingMaxHeaderListSize,
	}

	return TlsProfile{
		TransportSettings: TransportSettings{
			Spec:    nil,
			HelloID: tls.HelloChrome_140,
			Http2Settings: TransportHttp2Settings{
				Settings: settings,
				Order:    settingsOrder,
			},
		},
	}
}

type testClientTlsResponse struct {
	Ja4 string `json:"ja4"`
}

func TestClientTLSChrome140(t *testing.T) {
	expectedJa4 := "t13d1516h2_8daaf6152771_02713d6af862"

	client := MustNew(
		WithInsecureSkipVerify(),
		WithTlsProfile(chrome140Profile()),
	)

	res, err := client.Get(
		"https://tls.browserleaks.com/",
	)

	if err != nil {
		t.Fatalf("Error while sending request: %v", err)
	}

	var body testClientTlsResponse
	err = res.BodyDecode(&body)

	if err != nil {
		t.Fatalf("Error while binding response body: %v", err)
	}

	if body.Ja4 != expectedJa4 {
		t.Fatalf("Client has unexpected ja4, expected: %s got %s", expectedJa4, body.Ja4)
	}
}

func TestGetRequest(t *testing.T) {
	client := MustNew(
		WithInsecureSkipVerify(),
		WithTlsProfile(chrome140Profile()),
	)

	res, err := client.Get(
		fmt.Sprintf("http://127.0.0.1:%d/ping", testServerPort),
	)

	if err != nil {
		t.Fatalf("Unexpected error while pinging test server: %v", err)
	}

	if res.StatusCode() != 200 {
		t.Fatalf("Unexpected status code from test server: %s", res.Status())
	}

	if res.BodyString() != "pong" {
		t.Fatalf("Unexpected body from test server: %s", res.BodyString())
	}
}

func TestPostRequest(t *testing.T) {
	client := MustNew(
		WithInsecureSkipVerify(),
		WithTlsProfile(chrome140Profile()),
	)

	testCases := []struct {
		object   any
		expected string
	}{
		{
			object:   map[string]any{"foo": map[string]any{"bar": true}},
			expected: `{"foo":{"bar":true}}`,
		},
		{
			object: map[string]any{"foo": map[string]any{
				"bar":   true,
				"count": 1,
			}},
			expected: `{"foo":{"bar":true,"count":1}}`,
		},
	}

	for _, testCase := range testCases {
		res, err := client.Post(
			fmt.Sprintf("http://127.0.0.1:%d/json", testServerPort),
			RequestJsonBody(testCase.object),
		)

		if err != nil {
			t.Fatalf("Unexpected error while pinging test server: %v", err)
		}

		if res.StatusCode() != 200 {
			t.Fatalf("Unexpected status code from test server: %s", res.Status())
		}

		bodyString := res.BodyString()

		if bodyString != testCase.expected {
			t.Fatalf("Server responded with unexpected body, expected %s got %s", testCase.expected, bodyString)
		}
	}
}

func TestRequestCancelation(t *testing.T) {
	testCases := []struct {
		clientTimeoutDurationMiliseconds int
		serverTimeoutDurationMiliseconds int
		shouldError                      bool
	}{
		{clientTimeoutDurationMiliseconds: 100, serverTimeoutDurationMiliseconds: 500, shouldError: true},
		{clientTimeoutDurationMiliseconds: 120, serverTimeoutDurationMiliseconds: 500, shouldError: true},
		{clientTimeoutDurationMiliseconds: 150, serverTimeoutDurationMiliseconds: 500, shouldError: true},
		{clientTimeoutDurationMiliseconds: 200, serverTimeoutDurationMiliseconds: 500, shouldError: true},
		{clientTimeoutDurationMiliseconds: 50, serverTimeoutDurationMiliseconds: 500, shouldError: true},
		{clientTimeoutDurationMiliseconds: 100, serverTimeoutDurationMiliseconds: 500, shouldError: true},
		{clientTimeoutDurationMiliseconds: 300, serverTimeoutDurationMiliseconds: 500, shouldError: true},

		{clientTimeoutDurationMiliseconds: 1000, serverTimeoutDurationMiliseconds: 100, shouldError: false},
		{clientTimeoutDurationMiliseconds: 500, serverTimeoutDurationMiliseconds: 100, shouldError: false},
		{clientTimeoutDurationMiliseconds: 250, serverTimeoutDurationMiliseconds: 100, shouldError: false},
		{clientTimeoutDurationMiliseconds: 350, serverTimeoutDurationMiliseconds: 100, shouldError: false},
		{clientTimeoutDurationMiliseconds: 400, serverTimeoutDurationMiliseconds: 100, shouldError: false},
	}

	for _, testCase := range testCases {
		clientTimeoutDuration := time.Duration(testCase.clientTimeoutDurationMiliseconds) * time.Millisecond
		serverTimeoutDuration := time.Duration(testCase.serverTimeoutDurationMiliseconds) * time.Millisecond

		t.Logf(
			"Running testcase client timeout: %v server timeout: %v",
			clientTimeoutDuration,
			serverTimeoutDuration,
		)

		client := MustNew(
			WithCustomTimeout(clientTimeoutDuration),
			WithTlsProfile(chrome140Profile()),
		)

		res, err := client.Get(
			fmt.Sprintf("http://127.0.0.1:%d/timeout?timeoutMs=%d", testServerPort, testCase.serverTimeoutDurationMiliseconds),
		)

		if err != nil {
			if !errors.Is(err, ErrRequestTimedOut) {
				t.Fatalf("Unexpected error while getting timeout route: %v", err)
			}

			if !testCase.shouldError {
				t.Fatalf(
					"Request should pass, client duration: %v, server duration: %v err: %v",
					clientTimeoutDuration,
					serverTimeoutDuration,
					err,
				)
			}
		} else {
			if res.StatusCode() != http.StatusOK {
				t.Fatalf("Unexpected response status code: %s", res.Status())
			}
		}
	}
}
