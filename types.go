package http_client

import (
	"io"
	"time"

	fhttp "github.com/vimbing/fhttp"
	"github.com/vimbing/fhttp/cookiejar"
	"github.com/vimbing/fhttp/http2"
	tls "github.com/vimbing/utls"
)

type OptionStringJa string
type OptionTimeout time.Duration
type OptionProxy string
type OptionDisallowRedirect bool
type OptionForcedProxyRotation bool
type OptionTLSHelloID tls.ClientHelloID
type OptionTlsProfile TlsProfile
type OptionInsecureSkipVerify bool
type OptionCookieJar *cookiejar.Jar
type OptionRequestMiddleware []RequestMiddlewareFunc
type OptionResponseMiddleware []ResponseMiddlewareFunc
type OptionResponseErrorMiddleware []ResponseErrorMiddlewareFunc
type OptionRetry *Retry
type OptionStatusValidationFunc StatusValidationFunc

type Client struct {
	fhttpClient *fhttp.Client
	cfg         *Config
}

type RequestMiddlewareFunc func(*Request) error
type ResponseMiddlewareFunc func(*Response) error
type ResponseErrorMiddlewareFunc func(*Request, error)

type TransportHttp2Settings struct {
	Order    []http2.SettingID
	Settings map[http2.SettingID]uint32
}

type TransportSettings struct {
	Spec          *tls.ClientHelloSpec
	HelloID       tls.ClientHelloID
	Http2Settings TransportHttp2Settings
	Flow          uint32
}

type Config struct {
	insecureSkipVerify      bool
	requestMiddleware       []RequestMiddlewareFunc
	responseMiddleware      []ResponseMiddlewareFunc
	responseErrorMiddleware []ResponseErrorMiddlewareFunc
	proxies                 []string
	forceRotation           bool
	allowRedirect           bool
	timeout                 time.Duration
	jar                     *cookiejar.Jar
	transportSettings       TransportSettings
	retry                   *Retry
	statusValidationFunc    StatusValidationFunc
}

type RequestJsonBody any
type QueryParams map[string]string
type FormUrlEncoded map[string]string

type Request struct {
	Method string
	Body   io.Reader
	Header fhttp.Header
	Url    string

	// ctx       context.Context
	// ctxCancel context.CancelFunc

	protoMinor int
	protoMajor int
	proto      string

	host         *string
	fhttpRequest *fhttp.Request
}

type Response struct {
	Body          []byte
	fhttpResponse *fhttp.Response
}

type requestExecutionResult struct {
	res   *Response
	error error
}

type Retry struct {
	Max           int
	Delay         time.Duration
	IgnoredErrors []error
	EndingErrors  []error
	OnError       func(error)
}

type doFunc func(*Request) (*Response, error)

type StatusValidationFunc func(status int, client *Client) error

type TlsProfile struct {
	TransportSettings
}
