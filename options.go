package httpv3

import (
	"time"

	"github.com/repeale/fp-go"
	tls "github.com/vimbing/vutls"
)

type OptionStringJa string
type OptionTimeout time.Duration
type OptionProxy string
type OptionDisallowRedirect bool
type OptionUtlsJa3HelloId tls.ClientHelloID
type OptionUtlsJa3HelloSpec tls.ClientHelloSpec

func WithProxyListUnformatted(proxy string) []OptionProxy {
	return []OptionProxy{OptionProxy(proxy)}
}

func WithProxyList(proxy string) []OptionProxy {
	return []OptionProxy{OptionProxy(proxy)}
}

func WithProxyUnformatted(proxy string) OptionProxy {
	return OptionProxy(proxy)
}

func WithProxy(proxy string) OptionProxy {
	return OptionProxy(proxy)
}

func WithDisallowedRedirects() OptionDisallowRedirect {
	return false
}

func WithCustomTimeout(timeout time.Duration) OptionTimeout {
	return OptionTimeout(timeout)
}

func WithUtlsJa3Helloid(ja3HelloId tls.ClientHelloID) OptionUtlsJa3HelloId {
	return OptionUtlsJa3HelloId(ja3HelloId)
}

func parseOptions(options ...any) *Config {
	defaultCfg := &Config{
		proxies:       []string{},
		allowRedirect: true,
		timeout:       time.Second * 15,
		ja3:           tls.HelloChrome_120,
	}

	for _, opt := range options {
		switch v := opt.(type) {
		case OptionProxy:
			defaultCfg.proxies = []string{string(v)}
		case []OptionProxy:
			defaultCfg.proxies = fp.Map(func(p OptionProxy) string { return string(p) })(v)
		case OptionDisallowRedirect:
			defaultCfg.allowRedirect = false
		case OptionTimeout:
			defaultCfg.timeout = time.Duration(v)
		case OptionUtlsJa3HelloId:
			defaultCfg.ja3 = tls.ClientHelloID(v)
		}
	}

	return defaultCfg
}
