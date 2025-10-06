package http_client

import (
	"errors"
	"fmt"
	"os"
	"slices"
	"strings"
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
	Ja4        string `json:"ja4"`
	AkamaiHash string `json:"akamai_hash"`
	Ja3Text    string `json:"ja3_text"`
}

type ja3 struct {
	sslVersion               string
	ciphers                  []string
	sslExtensions            []string
	ellipticCurvePoints      []string
	ellipticCurvePointFormat string
}

func parseJa3Text(ja3Text string) *ja3 {
	ja3Parts := strings.Split(ja3Text, ",")

	sslVersion := ja3Parts[0]
	cipher := ja3Parts[1]
	sslExtension := ja3Parts[2]
	ellipticCurve := ja3Parts[3]
	ellipticCurvePointFormat := ja3Parts[4]

	return &ja3{
		sslVersion:               sslVersion,
		ciphers:                  strings.Split(cipher, "-"),
		sslExtensions:            strings.Split(sslExtension, "-"),
		ellipticCurvePoints:      strings.Split(ellipticCurve, "-"),
		ellipticCurvePointFormat: ellipticCurvePointFormat,
	}
}

func (j *ja3) compareWithExpected(expected *ja3, t *testing.T) {
	if j.sslVersion != expected.sslVersion {
		t.Errorf("Undexpected ssl version in ja3, expected: %s, got: %s", expected.sslVersion, j.sslVersion)
	}

	if j.ellipticCurvePointFormat != expected.ellipticCurvePointFormat {
		t.Errorf("Undexpected eliptic point format in ja3, expected: %s, got: %s", expected.ellipticCurvePointFormat, j.ellipticCurvePointFormat)
	}

	findAdditionalAndMissing := func(expected, actual []string) (missing []string, additional []string) {
		for _, cipher := range actual {
			if !slices.Contains(expected, cipher) {
				additional = append(additional, cipher)
			}
		}

		for _, cipher := range expected {
			if !slices.Contains(actual, cipher) {
				missing = append(missing, cipher)
			}
		}

		return
	}

	unexpectedMissingCiphers, unexpectedAdditionalCiphers := findAdditionalAndMissing(
		expected.ciphers,
		j.ciphers,
	)

	if len(unexpectedMissingCiphers) > 0 {
		t.Errorf("Found unexpected missing ciphers in ja3: [%s]", strings.Join(unexpectedMissingCiphers, ", "))
	}

	if len(unexpectedAdditionalCiphers) > 0 {
		t.Errorf("Found unexpected additional ciphers in ja3: [%s]", strings.Join(unexpectedAdditionalCiphers, ", "))
	}

	unexpectedMissingEllipticCurvePoints, unexpectedAdditionalEllipticCurvePoints := findAdditionalAndMissing(
		expected.ellipticCurvePoints,
		j.ellipticCurvePoints,
	)

	if len(unexpectedMissingEllipticCurvePoints) > 0 {
		t.Errorf("Found unexpected missing elliptic curves point in ja3: [%s]", strings.Join(unexpectedMissingEllipticCurvePoints, ", "))
	}

	if len(unexpectedAdditionalEllipticCurvePoints) > 0 {
		t.Errorf("Found unexpected additional elliptic curves point in ja3: [%s]", strings.Join(unexpectedAdditionalEllipticCurvePoints, ", "))
	}

	unexpectedMissingSSLExtensions, unexpectedAdditionalMissingSSLExtensions := findAdditionalAndMissing(
		expected.sslExtensions,
		j.sslExtensions,
	)

	if len(unexpectedMissingSSLExtensions) > 0 {
		t.Errorf("Found unexpected missing ssl extensions in ja3: [%s]", strings.Join(unexpectedMissingSSLExtensions, ", "))
	}

	if len(unexpectedAdditionalMissingSSLExtensions) > 0 {
		t.Errorf("Found unexpected additional ssl extensions in ja3: [%s]", strings.Join(unexpectedAdditionalMissingSSLExtensions, ", "))
	}

}

func TestClientTLSChrome140(t *testing.T) {
	// t13d1516h2_8daaf6152771_02713d6af862 maybe this as well
	expectedJa4 := "t13d1516h2_8daaf6152771_d8a2da3f94cd"
	expectedAkamaiHash := "52d84b11737d980aef856699f885ca86"
	expectedJaTextString := "771,4865-4866-4867-49195-49199-49196-49200-52393-52392-49171-49172-156-157-47-53,11-51-43-5-16-23-65281-0-45-35-65037-27-18-17613-13-10,4588-29-23-24,0"

	expectedJa3 := parseJa3Text(expectedJaTextString)

	client := MustNew(
		WithInsecureSkipVerify(),
		WithTlsProfile(chrome140Profile()),
	)

	res, err := client.Get(
		"https://tls.browserleaks.com/json",
		http.Header{http.PHeaderOrderKey: {":method", ":authority", ":scheme", ":path"}},
	)

	if err != nil {
		t.Fatalf("Error while sending request: %v", err)
	}

	var body testClientTlsResponse
	err = res.BodyDecode(&body)

	if err != nil {
		t.Fatalf("Error while binding response body: %v", err)
	}

	clientJa3 := parseJa3Text(body.Ja3Text)
	clientJa3.compareWithExpected(expectedJa3, t)

	if body.Ja4 != expectedJa4 {
		t.Errorf("Client has unexpected ja4, expected: %s got %s", expectedJa4, body.Ja4)
	}

	if body.AkamaiHash != expectedAkamaiHash {
		t.Errorf("Client has unexpected akamai hash, expected: %s got %s", expectedAkamaiHash, body.AkamaiHash)
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
		t.Errorf("Unexpected body from test server: %s", res.BodyString())
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
			t.Errorf("Server responded with unexpected body, expected %s got %s", testCase.expected, bodyString)
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

func TestProxyRotation(t *testing.T) {
	proxy := os.Getenv("TEST_PROXY")

	if len(proxy) == 0 {
		t.Skip("Omiting proxy test, becouse no proxy was provided in env. If you want to perform this test, set env variable: TEST_PROXY=host:port@user:password;host2:port2@user2:password2")
		return
	}

	proxies := strings.Split(proxy, ";")
	lastAddress := ""

	client := MustNew(
		WithTlsProfile(chrome140Profile()),
	)

	for _, proxy := range proxies {
		client.ChangeProxy(proxy)

		res, err := client.Get("https://icanhazip.com")

		if err != nil {
			t.Fatalf("Unexpected error while getting ip service: %v", err)
		}

		t.Logf("Ip service responded with ip: %s", res.BodyString())

		if lastAddress == res.BodyString() {
			t.Errorf("Proxy repeated after rotation: %s", lastAddress)
		}
	}
}
