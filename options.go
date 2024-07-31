package http

import (
	"time"

	"github.com/repeale/fp-go"
	lo "github.com/samber/lo"
	"github.com/vimbing/fhttp/cookiejar"
	tls "github.com/vimbing/vutls"
)

func WithProxyList(proxyList []string) []OptionProxy {
	return append([]OptionProxy{}, parseList(proxyList)...)
}

func WithProxyListParsed(proxyList []string) []OptionProxy {
	return append([]OptionProxy{}, lo.Map(proxyList, func(p string, i int) OptionProxy { return OptionProxy(p) })...)
}

func WithProxy(proxy string) OptionProxy {
	parsed, _ := parseSingleProxy(proxy)
	return OptionProxy(parsed)
}

func WithProxyParsed(proxy string) OptionProxy {
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

func WithTlsProfile(tlsProfile TlsProfile) OptionTlsProfile {
	return OptionTlsProfile(tlsProfile)
}

func WithInsecureSkipVerify() OptionInsecureSkipVerify {
	return OptionInsecureSkipVerify(true)
}

func WithCookieJar(jar *cookiejar.Jar) OptionCookieJar {
	return OptionCookieJar(jar)
}

func parseOptions(options ...any) *Config {
	defaultCfg := &Config{
		proxies:       []string{},
		allowRedirect: true,
		timeout:       time.Second * 15,
		ja3:           tls.HelloChrome_120,
		jar:           nil,
	}

	for _, opt := range options {
		switch v := opt.(type) {
		case OptionProxy:
			defaultCfg.proxies = []string{string(v)}
		case []OptionProxy:
			defaultCfg.proxies = fp.Map(func(p OptionProxy) string { return string(p) })(v)
		case OptionDisallowRedirect:
			defaultCfg.allowRedirect = false
		case OptionCookieJar:
			defaultCfg.jar = v
		case OptionTimeout:
			defaultCfg.timeout = time.Duration(v)
		case OptionUtlsJa3HelloId:
			defaultCfg.ja3 = tls.ClientHelloID(v)
		case OptionTlsProfile:
			profile := TlsProfile(v)
			defaultCfg.tlsProfile = &profile
		case OptionInsecureSkipVerify:
			defaultCfg.insecureSkipVerify = true
		}
	}

	return defaultCfg
}
