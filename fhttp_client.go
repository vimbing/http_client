package http_client

import (
	"errors"
	"net"
	"net/url"
	"strings"
	"time"

	http "github.com/vimbing/fhttp"
	"github.com/vimbing/retry"

	"golang.org/x/net/proxy"
)

func socksDialer(pickedProxy string) (proxy.ContextDialer, error) {
	proxyUrl, err := url.Parse(pickedProxy)

	if err != nil {
		return nil, err
	}

	var auth *proxy.Auth

	if proxyUrl.User != nil && len(proxyUrl.User.Username()) > 0 {
		password, passwordSet := proxyUrl.User.Password()

		auth = &proxy.Auth{
			User: proxyUrl.User.Username(),
		}

		if passwordSet {
			auth.Password = password
		}
	}

	dialSocksProxy, err := proxy.SOCKS5("tcp", proxyUrl.Host, auth, &net.Dialer{
		Timeout:   30 * time.Second,
		KeepAlive: 30 * time.Second,
	})

	if err != nil {
		return nil, err
	}

	var ok bool

	dialer, ok := dialSocksProxy.(proxy.ContextDialer)

	if !ok {
		return nil, errors.New("failed type assertion to DialContext")
	}

	return dialer, nil
}

// rebindRoundtripper rebinds the Transport of the provided http.Client using the
// configuration in cfg by selecting and configuring an appropriate proxy dialer
// and creating a new round tripper with the transport settings from cfg.
//
// If cfg.proxies contains entries, a proxy is chosen at random; a SOCKS5 dialer
// is used when the proxy string contains "socks", otherwise a CONNECT dialer is
// created. When no proxies are configured, proxy.Direct is used. The new
// Transport is configured with cfg.transportSettings.HelloID, cfg.insecureSkipVerify,
// the chosen dialer, and HTTP/2 settings from cfg.transportSettings.Http2Settings.
//
// The operation is retried up to three times on error. It returns any error
// encountered while creating the dialer or configuring the Transport.
func rebindRoundtripper(c *http.Client, cfg *Config) error {
	return retry.Retrier{Max: 3, Delay: time.Second * 0}.Retry(func() error {
		var dialer proxy.ContextDialer
		var err error

		if len(cfg.proxies) > 0 {
			pickedProxy := cfg.proxies[RandomInt(0, len(cfg.proxies))]

			if strings.Contains(cfg.proxies[0], "socks") {
				dialer, err = socksDialer(pickedProxy)
			} else {
				dialer, err = newConnectDialer(pickedProxy)
			}
		} else {
			dialer = proxy.Direct
		}

		if err != nil {
			return err
		}

		c.Transport = newRoundTripper(roundTripperSettings{
			clientHello:        cfg.transportSettings.HelloID,
			insecureSkipVerify: cfg.insecureSkipVerify,
			dialer:             dialer,
			http2Settings:      cfg.transportSettings.Http2Settings.Settings,
			http2SettingsOrder: cfg.transportSettings.Http2Settings.Order,
		})

		return nil
	})
}

func newFhttpClient(cfg *Config) (*http.Client, error) {
	client := &http.Client{
		Timeout: cfg.timeout,
	}

	if !cfg.allowRedirect {
		client.CheckRedirect = func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		}
	}

	client.Timeout = cfg.timeout

	if err := rebindRoundtripper(client, cfg); err != nil {
		return client, err
	}

	return client, nil
}