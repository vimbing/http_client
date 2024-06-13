package httpv3

import (
	"io"
	"time"

	fhttp "github.com/vimbing/fhttp"
	tls "github.com/vimbing/vutls"
)

type Client struct {
	fhttpClient *fhttp.Client
	cfg         *Config
}

type RequestMiddlewareFunc func(*Request) error
type ResponseMiddlewareFunc func(*Response) error

type Config struct {
	requestMiddleware  []RequestMiddlewareFunc
	responseMiddleware []ResponseMiddlewareFunc
	proxies            []string
	allowRedirect      bool
	timeout            time.Duration
	ja3                tls.ClientHelloID
}

type RequestJsonBody any

type Request struct {
	Method string
	Body   io.Reader
	Header fhttp.Header
	Url    string

	fhttpRequest *fhttp.Request
}

type Response struct {
	body          []byte
	fhttpResponse *fhttp.Response
}
