package httpv3

import fhttp "github.com/vimbing/fhttp"

func (r *Request) Build() error {
	req, err := fhttp.NewRequest(r.Method, r.Url, r.Body)

	if err != nil {
		return err
	}

	req.Header = r.Header
	r.fhttpRequest = req

	return nil
}
