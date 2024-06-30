package http

import (
	"context"
	"strings"
	"time"

	fhttp "github.com/vimbing/fhttp"
)

func (r *Request) useTlsProfile() {
	if r.tlsProfile == nil {
		return
	}

	for k := range r.Header {
		switch strings.ToLower(k) {
		case "sec-ch-ua":
			r.Header[k] = []string{r.tlsProfile.SecChUa}
		case "sec-ch-ua-mobile":
			r.Header[k] = []string{r.tlsProfile.SecChUaMobile}
		case "sec-ch-ua-platform":
			r.Header[k] = []string{r.tlsProfile.SecChaUaPlatform}
		case "user-agent":
			r.Header[k] = []string{r.tlsProfile.UserAgent}
		}
	}
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

	return cancel, nil
}
