package http

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
			clientHello:        cfg.ja3,
			insecureSkipVerify: cfg.insecureSkipVerify,
			dialer:             dialer,
			http2Settings:      cfg.httpSettings.Settings,
			http2SettingsOrder: cfg.httpSettings.Order,
			disablePush:        cfg.httpSettings.DisablePush,
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
