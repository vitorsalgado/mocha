package mocha

import (
	"crypto/tls"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"time"

	"github.com/rs/zerolog"

	"github.com/vitorsalgado/mocha/v3/misc"
)

var _ ProxyConfigurer = (*ProxyConfig)(nil)

// ProxyConfig configures a proxy.
type ProxyConfig struct {
	// Via sets a URL to route the request via another proxy server.
	// Via is only valid when Transport configuration is not set.
	Via string

	// Timeout is the timeout used when calling the proxy client.
	Timeout time.Duration

	// SSLVerify enable/disable server certificate verification.
	SSLVerify bool

	// Transport sets a custom http.RoundTripper.
	// Target config will be ignored. Set it manually in your http.RoundTripper implementation.
	// If none is provided, a default one will be used.
	// You need to set a Target and a custom Via in your custom http.RoundTripper.
	Transport http.RoundTripper
}

// ProxyConfigurer lets users configure the proxy.
type ProxyConfigurer interface {
	Apply(config *ProxyConfig) error
}

// Apply allows ProxyConfig to be used as a Configurer.
func (p *ProxyConfig) Apply(c *ProxyConfig) error {
	c.Transport = p.Transport
	c.Timeout = p.Timeout
	c.Via = p.Via
	c.SSLVerify = p.SSLVerify

	return nil
}

var _defaultProxyConfig = ProxyConfig{Timeout: 10 * time.Second, SSLVerify: false}

type reverseProxy struct {
	log          *zerolog.Logger
	roundTripper http.RoundTripper
}

func newProxy(log *zerolog.Logger, config *Config) *reverseProxy {
	var roundTripper http.RoundTripper

	if config.Proxy.Transport == nil {
		transport := &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: !config.Proxy.SSLVerify, ClientCAs: config.TLSClientCAs},
		}

		if config.TLSCertificate != nil {
			transport.TLSClientConfig.Certificates = []tls.Certificate{*config.TLSCertificate}
		}

		if config.Proxy.Via != "" {
			u, err := url.Parse(config.Proxy.Via)
			if err != nil {
				panic(fmt.Errorf("proxy: failed to parse proxy via url %s: %w", config.Proxy.Via, err))
			}

			transport.Proxy = http.ProxyURL(u)
		}

		transport.TLSHandshakeTimeout = 15 * time.Second
		transport.IdleConnTimeout = 15 * time.Second
		transport.ExpectContinueTimeout = 1 * time.Second
		transport.ResponseHeaderTimeout = config.Proxy.Timeout

		roundTripper = transport
	} else {
		roundTripper = config.Proxy.Transport
	}

	return &reverseProxy{log, roundTripper}
}

func (p *reverseProxy) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodConnect {
		p.handleTunneling(w, r)
	} else {
		p.handleHTTP(w, r)
	}
}

func (p *reverseProxy) handleTunneling(w http.ResponseWriter, r *http.Request) {
	hijacker, ok := w.(http.Hijacker)
	if !ok {
		http.Error(w, "proxy: hijacking is not supported", http.StatusInternalServerError)
		return
	}

	in, _, err := hijacker.Hijack()
	if err != nil {
		http.Error(w, err.Error(), http.StatusServiceUnavailable)
	}

	defer in.Close()

	out, err := net.DialTimeout("tcp", r.Host, 5*time.Second)
	if err != nil {
		http.Error(w, err.Error(), http.StatusServiceUnavailable)
		return
	}

	defer out.Close()

	w.WriteHeader(http.StatusOK)

	errCh := make(chan error, 2)
	cp := func(dst io.WriteCloser, src io.ReadCloser) {
		_, err = io.Copy(dst, src)
		errCh <- err
	}

	go cp(out, in)
	go cp(in, out)

	err = <-errCh
	if err != nil {
		p.log.Error().Err(err).
			Str("url", r.URL.String()).
			Str("method", r.Method).
			Msg("proxy: error writing response")
	}
}

func (p *reverseProxy) handleHTTP(w http.ResponseWriter, r *http.Request) {
	r.Header.Del(misc.HeaderConnection)

	res, err := p.roundTripper.RoundTrip(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusServiceUnavailable)
		return
	}

	defer res.Body.Close()

	for k, vv := range res.Header {
		for _, v := range vv {
			w.Header().Add(k, v)
		}
	}

	w.WriteHeader(res.StatusCode)

	_, err = io.Copy(w, res.Body)
	if err != nil {
		p.log.Error().Err(err).
			Str("url", r.URL.String()).
			Str("method", r.Method).
			Str("status", res.Status).
			Msg("proxy: error copying response body")
	}

	for k, vv := range res.Trailer {
		for _, v := range vv {
			w.Header().Add(k, v)
		}
	}
}
