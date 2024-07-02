package http

import (
	"bytes"
	"compress/gzip"
	"compress/zlib"
	"io"
	"strings"

	"github.com/andybalholm/brotli"
	"github.com/klauspost/compress/zstd"
	http "github.com/vimbing/fhttp"
)

func decodeResponseBody(headers http.Header, body []byte) ([]byte, error) {
	var encoding string

	for k, v := range headers {
		if strings.EqualFold(k, "content-encoding") && len(v) > 0 {
			encoding = strings.ToLower(v[0])
		}
	}

	if len(encoding) == 0 {
		return body, nil
	}

	switch encoding {
	case "br":
		return brotoliDecode(body)
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

func brotoliDecode(data []byte) ([]byte, error) {
	return io.ReadAll(brotli.NewReader(bytes.NewReader(data)))
}

func deflateDecode(data []byte) ([]byte, error) {
	zr, err := zlib.NewReader(bytes.NewReader(data))

	if err != nil {
		return []byte{}, err
	}

	defer zr.Close()

	return io.ReadAll(zr)
}

func gzipDecode(data []byte) ([]byte, error) {
	gz, err := gzip.NewReader(bytes.NewReader(data))

	if err != nil {
		return []byte{}, err
	}

	defer gz.Close()

	return io.ReadAll(gz)
}

func zstdDecode(data []byte) ([]byte, error) {
	zs, err := zstd.NewReader(bytes.NewReader(data))

	if err != nil {
		return []byte{}, err
	}

	defer zs.Close()

	return io.ReadAll(zs)
}
