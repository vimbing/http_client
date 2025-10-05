package http_client

import (
	"compress/gzip"
	"compress/zlib"
	"io"
	"strings"

	"github.com/andybalholm/brotli"
	"github.com/klauspost/compress/zstd"
	http "github.com/vimbing/fhttp"
)

func decodeResponseBody(headers http.Header, body io.Reader) (io.Reader, error) {
	var encoding string

	for k, v := range headers {
		if strings.EqualFold(k, "content-encoding") && len(v) > 0 {
			encoding = strings.ToLower(v[0])
		}
	}

	switch encoding {
	case "br":
		return brotoliDecode(body), nil
	case "gzip":
		return gzipDecode(body)
	case "deflate":
		return deflateDecode(body)
	case "zstd":
		return zstdDecode(body)
	default:
		return body, nil
	}
}

func brotoliDecode(body io.Reader) io.Reader {
	return brotli.NewReader(body)
}

func deflateDecode(body io.Reader) (io.Reader, error) {
	zr, err := zlib.NewReader(body)
	return zr, err
}

func gzipDecode(body io.Reader) (io.Reader, error) {
	gz, err := gzip.NewReader(body)
	return gz, err
}

func zstdDecode(body io.Reader) (io.Reader, error) {
	zs, err := zstd.NewReader(body)
	return zs, err
}
