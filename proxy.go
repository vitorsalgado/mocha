package mocha

import (
	"io"
	"net"
	"net/http"
	"net/url"
	"time"

	"github.com/vitorsalgado/mocha/v3/event"
	"github.com/vitorsalgado/mocha/v3/internal/header"
)

// ProxyConfig configures proxy.
type ProxyConfig struct {
	// Target configures the proxy address to route requests
	// by the actual proxy being set.
	// You need to set this configuration in your custom Transport, if one is provided.
	Target string

	// Timeout is the timeout used when calling the proxy client.
	Timeout time.Duration

	// Transport sets a custom http.RoundTripper.
	// Target config will be ignored. Set it manually in your http.RoundTripper implementation.
	// If none is provided, a default one will be used.
	Transport http.RoundTripper
}

// ProxyConfigurer lets users configure proxy.
type ProxyConfigurer interface {
	Apply(config *ProxyConfig)
}

// Apply allows ProxyConfig to be used as a Configurer.
func (p *ProxyConfig) Apply(c *ProxyConfig) {
	c.Target = p.Target
	c.Transport = p.Transport
	c.Timeout = p.Timeout
}

var _defaultProxyConfig = ProxyConfig{Timeout: 10 * time.Second}

type proxy struct {
	conf     *ProxyConfig
	listener *event.Listener
}

func newProxy(conf *ProxyConfig, events *event.Listener) *proxy {
	if conf.Transport == nil {
		transport := &http.Transport{}
		if conf.Target != "" {
			u, err := url.Parse(conf.Target)
			if err != nil {
				panic(err)
			}

			transport.Proxy = http.ProxyURL(u)
		}

		transport.TLSHandshakeTimeout = 15 * time.Second
		transport.IdleConnTimeout = 15 * time.Second
		transport.ExpectContinueTimeout = 1 * time.Second
		transport.ResponseHeaderTimeout = conf.Timeout

		conf.Transport = transport
	}

	return &proxy{conf: conf, listener: events}
}

func (p *proxy) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodConnect {
		p.handleTunneling(w, r)
	} else {
		p.handleHTTP(w, r)
	}
}

func (p *proxy) handleTunneling(w http.ResponseWriter, r *http.Request) {
	hijacker, ok := w.(http.Hijacker)
	if !ok {
		http.Error(w, "[proxy] hijacking is not supported", http.StatusInternalServerError)
		return
	}

	in, _, err := hijacker.Hijack()
	if err != nil {
		p.listener.Emit(&event.OnError{Request: event.FromRequest(r), Err: err})
		http.Error(w, err.Error(), http.StatusServiceUnavailable)
	}

	defer in.Close()

	out, err := net.DialTimeout("tcp", r.Host, 5*time.Second)
	if err != nil {
		p.listener.Emit(&event.OnError{Request: event.FromRequest(r), Err: err})
		http.Error(w, err.Error(), http.StatusServiceUnavailable)
		return
	}

	defer out.Close()

	w.WriteHeader(http.StatusOK)

	errCh := make(chan error, 2)
	cp := func(dst io.WriteCloser, src io.ReadCloser) {
		_, err := io.Copy(dst, src)
		errCh <- err
	}

	go cp(out, in)
	go cp(in, out)

	err = <-errCh
	if err != nil {
		p.listener.Emit(&event.OnError{Request: event.FromRequest(r), Err: err})
	}
}

func (p *proxy) handleHTTP(w http.ResponseWriter, r *http.Request) {
	r.Header.Del(header.Connection)

	res, err := p.conf.Transport.RoundTrip(r)
	if err != nil {
		p.listener.Emit(&event.OnError{Request: event.FromRequest(r), Err: err})
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
		p.listener.Emit(&event.OnError{Request: event.FromRequest(r), Err: err})
	}
}
