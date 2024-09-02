package profiles

import (
	fhttp "github.com/vimbing/fhttp"
	"github.com/vimbing/fhttp/http2"
	http "github.com/vimbing/httpv2"
)

var HelloChrome_120 = http.TlsProfile{
	Headers: fhttp.Header{
		"sec-ch-ua":          {`"Not/A)Brand";v="8", "Chromium";v="126", "Google Chrome";v="126"`},
		"sec-ch-ua-mobile":   {"?0"},
		"sec-ch-ua-platform": {"\"Windows\""},
		"user-agent":         {"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/126.0.0.0 Safari/537.36"},
	},
}

var HelloZalandoIOS = http.TlsProfile{
	Headers: fhttp.Header{
		"user-agent": {`Mozilla/5.0 (iPhone; CPU iPhone OS 16_0 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/16.0 Mobile/15E148 Safari/604.1`},
	},
	Http2Settings: http.Http2Settings{
		Settings: map[http2.SettingID]uint32{
			http2.SettingHeaderTableSize:      4096,
			http2.SettingMaxConcurrentStreams: 100,
			http2.SettingInitialWindowSize:    2097152,
			http2.SettingMaxFrameSize:         16384,
			http2.SettingMaxHeaderListSize:    4294967295,
		},
		Order: []http2.SettingID{
			http2.SettingHeaderTableSize,
			http2.SettingMaxConcurrentStreams,
			http2.SettingInitialWindowSize,
			http2.SettingMaxFrameSize,
			http2.SettingMaxHeaderListSize,
		},
		DisablePush: true,
	},
	PseudoHeaderOrder: []string{
		":method",
		":path",
		":authority",
		":scheme",
	},
}
