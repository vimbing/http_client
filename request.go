package http

import (
	"context"
	"time"

	fhttp "github.com/vimbing/fhttp"
)

func (r *Request) useTlsProfile() {
	if r.tlsProfile == nil {
		return
	}

	r.Header.Set("sec-ch-ua", r.tlsProfile.SecChUa)
	r.Header.Set("sec-ch-ua-mobile", r.tlsProfile.SecChUaMobile)
	r.Header.Set("sec-ch-ua-platform", r.tlsProfile.SecChaUaPlatform)
	r.Header.Set("user-agent", r.tlsProfile.UserAgent)
}

func (r *Request) Build(timeout time.Duration) (context.CancelFunc, error) {
	ctx, cancel := context.WithTimeout(
		context.Background(),
		timeout,
	)

	r.useTlsProfile()

	req, err := fhttp.NewRequestWithContext(ctx, r.Method, r.Url, r.Body)

	if err != nil {
		return cancel, err
	}

	req.Header = r.Header
	r.fhttpRequest = req

	if r.host != nil {
		r.fhttpRequest.Host = *r.host
	}

	return cancel, nil
}

func (r *Request) SetHost(host string) {
	r.host = &host
}
