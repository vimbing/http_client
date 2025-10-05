package http_client

import (
	"time"

	"github.com/repeale/fp-go"
	lo "github.com/samber/lo"
	"github.com/vimbing/fhttp/cookiejar"
	tls "github.com/vimbing/utls"
)

func WithForcedProxyRotation() OptionForcedProxyRotation {
	return true
}

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

func WithProxySocks(proxy string) OptionProxy {
	parsed, _ := parseSingleSocksProxy(proxy)
	return OptionProxy(parsed)
}

func WithProxyParsed(proxy string) OptionProxy {
	return OptionProxy(proxy)
}

func WithTlsProfile(profile TlsProfile) OptionTlsProfile {
	return OptionTlsProfile(profile)
}

func WithDisallowedRedirects() OptionDisallowRedirect {
	return false
}

func WithCustomTimeout(timeout time.Duration) OptionTimeout {
	return OptionTimeout(timeout)
}

func WithInsecureSkipVerify() OptionInsecureSkipVerify {
	return OptionInsecureSkipVerify(true)
}

func WithCookieJar(jar *cookiejar.Jar) OptionCookieJar {
	return OptionCookieJar(jar)
}

func WithRequestMiddleware(m ...RequestMiddlewareFunc) OptionRequestMiddleware {
	return OptionRequestMiddleware(m)
}

func WithResponseMiddleware(m ...ResponseMiddlewareFunc) OptionResponseMiddleware {
	return OptionResponseMiddleware(m)
}

func WithResponseErrorMiddleware(m ...ResponseErrorMiddlewareFunc) OptionResponseErrorMiddleware {
	return OptionResponseErrorMiddleware(m)
}

func WithRetry(retry *Retry) OptionRetry {
	return OptionRetry(retry)
}

func WithStatusValidation(f StatusValidationFunc) OptionStatusValidationFunc {
	return OptionStatusValidationFunc(f)
}

func parseOptions(options ...any) *Config {
	defaultCfg := &Config{
		proxies:              []string{},
		allowRedirect:        true,
		timeout:              time.Second * 15,
		transportSettings:    TransportSettings{},
		jar:                  nil,
		retry:                &Retry{},
		statusValidationFunc: nil,
	}

	for _, opt := range options {
		switch v := opt.(type) {
		case OptionForcedProxyRotation:
			defaultCfg.forceRotation = true
		case OptionProxy:
			defaultCfg.proxies = []string{string(v)}
		case []OptionProxy:
			defaultCfg.proxies = fp.Map(func(p OptionProxy) string { return string(p) })(v)
		case OptionDisallowRedirect:
			defaultCfg.allowRedirect = false
		case OptionCookieJar:
			defaultCfg.jar = v
		case OptionStatusValidationFunc:
			defaultCfg.statusValidationFunc = StatusValidationFunc(v)
		case OptionResponseMiddleware:
			for _, m := range v {
				defaultCfg.responseMiddleware = append(defaultCfg.responseMiddleware, ResponseMiddlewareFunc(m))
			}
		case OptionRequestMiddleware:
			for _, m := range v {
				defaultCfg.requestMiddleware = append(defaultCfg.requestMiddleware, RequestMiddlewareFunc(m))
			}
		case OptionResponseErrorMiddleware:
			for _, m := range v {
				defaultCfg.responseErrorMiddleware = append(defaultCfg.responseErrorMiddleware, ResponseErrorMiddlewareFunc(m))
			}
		case OptionTimeout:
			defaultCfg.timeout = time.Duration(v)
		case OptionTLSHelloID:
			defaultCfg.transportSettings.HelloID = tls.ClientHelloID(v)
		case OptionTlsProfile:
			p := TlsProfile(v)
			defaultCfg.transportSettings = p.TransportSettings

			// bogdanHelloID := p.GetClientHelloId()
			// bogdanSpec, err := p.GetClientHelloSpec()
			// var spec *tls.ClientHelloSpec

			// if err == nil {
			// 	spec = &tls.ClientHelloSpec{
			// 		CipherSuites:       bogdanSpec.CipherSuites,
			// 		CompressionMethods: bogdanSpec.CompressionMethods,
			// 	}

			// 	for _, extension := range bogdanSpec.Extensions {
			// 		spec.Extensions = append(spec.Extensions, extension.(tls.TLSExtension))
			// 	}
			// }

			// helloID := tls.ClientHelloID{
			// 	Client:  bogdanHelloID.Client,
			// 	Version: bogdanHelloID.Version,
			// 	Seed:    (*tls.PRNGSeed)(bogdanHelloID.Seed),
			// 	Weights: &tls.DefaultWeights,
			// }

			// defaultCfg.transportSettings.helloID = helloID
			// defaultCfg.transportSettings.Flow = p.GetConnectionFlow()

			// http2Settings := map[http2.SettingID]uint32{}
			// http2SettingsOrder := []http2.SettingID{}

			// for _, settingOrderSettingID := range p.GetSettingsOrder() {
			// 	http2SettingsOrder = append(http2SettingsOrder, http2.SettingID(settingOrderSettingID))
			// }

			// for key, value := range p.GetSettings() {
			// 	if http2.SettingEnablePush == http2.SettingID(key) {
			// 		if (value) == 0 {
			// 			defaultCfg.transportSettings.DisablePush = true
			// 		}
			// 	}

			// 	http2Settings[http2.SettingID(key)] = value
			// }

			// defaultCfg.transportSettings.Http2Settings = TransportHttp2Settings{
			// 	Order:    http2SettingsOrder,
			// 	Settings: http2Settings,
			// }

			// TODO
			// defaultCfg.transportSettings.DisablePush = p
		case OptionInsecureSkipVerify:
			defaultCfg.insecureSkipVerify = true
		case OptionRetry:
			defaultCfg.retry = v
		}
	}

	return defaultCfg
}
