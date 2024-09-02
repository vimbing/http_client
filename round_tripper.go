package http

import (
	"context"
	"errors"
	"fmt"
	"net"
	"strings"
	"sync"

	http "github.com/vimbing/fhttp"

	"github.com/vimbing/fhttp/http2"
	"golang.org/x/net/proxy"

	utls "github.com/vimbing/vutls"
)

var errProtocolNegotiated = errors.New("protocol negotiated")

type roundTripper struct {
	sync.Mutex

	insecureSkipVerify bool
	clientHelloId      utls.ClientHelloID

	cachedConnections map[string]net.Conn
	cachedTransports  map[string]http.RoundTripper

	dialer proxy.ContextDialer

	http2Settings      map[http2.SettingID]uint32
	http2SettingsOrder []http2.SettingID
	disablePush        bool
}

func (rt *roundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	addr := rt.getDialTLSAddr(req)

	rt.Lock()
	if _, ok := rt.cachedTransports[addr]; !ok {
		rt.Unlock()
		if err := rt.getTransport(req, addr); err != nil {
			return nil, err
		}
	} else {
		rt.Unlock()
	}

	rt.Lock()
	transport := rt.cachedTransports[addr]
	rt.Unlock()
	return transport.RoundTrip(req)
}

func (rt *roundTripper) getTransport(req *http.Request, addr string) error {
	switch strings.ToLower(req.URL.Scheme) {
	case "http":
		rt.Lock()
		rt.cachedTransports[addr] = &http.Transport{DialContext: rt.dialer.DialContext}
		rt.Unlock()
		return nil
	case "https":
	default:
		return fmt.Errorf("invalid URL scheme: [%v]", req.URL.Scheme)
	}

	_, err := rt.dialTLS(context.Background(), "tcp", addr)
	switch err {
	case errProtocolNegotiated:
	case nil:
		// Should never happen.
		return errors.New("dialTLS returned no error when determining cachedTransports")
	default:
		return err
	}

	return nil
}

func (rt *roundTripper) dialTLS(ctx context.Context, network, addr string) (net.Conn, error) {
	rt.Lock()
	defer rt.Unlock()

	// If we have the connection from when we determined the HTTPS
	// cachedTransports to use, return that.
	if conn := rt.cachedConnections[addr]; conn != nil {
		delete(rt.cachedConnections, addr)
		return conn, nil
	}

	rawConn, err := rt.dialer.DialContext(ctx, network, addr)
	if err != nil {
		return nil, err
	}

	var host string
	if host, _, err = net.SplitHostPort(addr); err != nil {
		host = addr
	}

	conn := utls.UClient(rawConn, &utls.Config{
		ServerName:         host,
		InsecureSkipVerify: rt.insecureSkipVerify,
	},
		rt.clientHelloId,
	)

	if err != nil {
		_ = conn.Close()
		return nil, err
	}

	if err = conn.Handshake(); err != nil {
		_ = conn.Close()
		return nil, err
	}

	if rt.cachedTransports[addr] != nil {
		return conn, nil
	}

	// No http.Transport constructed yet, create one based on the results
	// of ALPN.
	switch conn.ConnectionState().NegotiatedProtocol {

	case http2.NextProtoTLS:
		t2 := http2.Transport{DialTLS: rt.dialTLSHTTP2}

		if len(rt.http2Settings) == 0 || len(rt.http2SettingsOrder) == 0 {
			t2.HeaderTableSize = 65536

			t2.Settings = []http2.Setting{
				{ID: http2.SettingMaxConcurrentStreams, Val: 1000},
				{ID: http2.SettingMaxHeaderListSize, Val: 262144},
			}

			t2.InitialWindowSize = 6291456
		} else {
			for _, settingId := range rt.http2SettingsOrder {
				t2.Settings = append(t2.Settings, http2.Setting{
					ID:  settingId,
					Val: rt.http2Settings[settingId],
				})
			}
		}

		if !rt.disablePush {
			t2.PushHandler = &http2.DefaultPushHandler{}
		}

		rt.cachedTransports[addr] = &t2
	default:
		// Assume the remote peer is speaking HTTP 1.x + TLS.
		rt.cachedTransports[addr] = &http.Transport{DialTLSContext: rt.dialTLS}
	}

	// Stash the connection just established for use servicing the
	// actual request (should be near-immediate).
	rt.cachedConnections[addr] = conn

	return nil, errProtocolNegotiated
}

func (rt *roundTripper) dialTLSHTTP2(network, addr string, _ *utls.Config) (net.Conn, error) {
	return rt.dialTLS(context.Background(), network, addr)
}

func (rt *roundTripper) getDialTLSAddr(req *http.Request) string {
	host, port, err := net.SplitHostPort(req.URL.Host)
	if err == nil {
		return net.JoinHostPort(host, port)
	}
	return net.JoinHostPort(req.URL.Host, "443") // we can assume port is 443 at this point
}

type roundTripperSettings struct {
	clientHello        utls.ClientHelloID
	insecureSkipVerify bool
	dialer             proxy.ContextDialer
	http2Settings      map[http2.SettingID]uint32
	http2SettingsOrder []http2.SettingID
	disablePush        bool
}

func newRoundTripper(settings roundTripperSettings) http.RoundTripper {
	return &roundTripper{
		dialer:             settings.dialer,
		insecureSkipVerify: settings.insecureSkipVerify,
		clientHelloId:      settings.clientHello,
		cachedTransports:   make(map[string]http.RoundTripper),
		cachedConnections:  make(map[string]net.Conn),
		http2Settings:      settings.http2Settings,
		http2SettingsOrder: settings.http2SettingsOrder,
		disablePush:        settings.disablePush,
	}
}
