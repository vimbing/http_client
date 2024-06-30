package profiles

import http "github.com/vimbing/httpv2"

var HelloChrome_120 = http.TlsProfile{
	SecChUa:          `"Not/A)Brand";v="8", "Chromium";v="126", "Google Chrome";v="126"`,
	SecChUaMobile:    "?0",
	SecChaUaPlatform: "\"Windows\"",
	UserAgent:        "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/126.0.0.0 Safari/537.36",
}
