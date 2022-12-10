package mocha

import (
	"io"
	"net"
	"net/http"
	"net/url"
	"time"
)

// ProxyConfig configures proxy.
type ProxyConfig struct {
	// ProxyVia configures the proxy address to route requests
	// by the actual proxy being set.
	// You need to set this configuration in your custom Transport, if one is provided.
	ProxyVia *url.URL

	// Transport sets a custom http.RoundTripper.
	// ProxyVia config will be ignored. Set it manually in your http.RoundTripper implementation.
	// If none is provided, a default one will be used.
	Transport http.RoundTripper
}

// ProxyConfigurer lets users configure proxy.
type ProxyConfigurer interface {
	Apply(config *ProxyConfig)
}

// Apply allows ProxyConfig to be used as a Configurer.
func (p *ProxyConfig) Apply(c *ProxyConfig) {
	c.ProxyVia = p.ProxyVia
	c.Transport = p.Transport
}

type proxy struct {
	conf *ProxyConfig
	e    *eventListener
}

func newProxy(conf *ProxyConfig, events *eventListener) *proxy {
	if conf.Transport == nil {
		transport := &http.Transport{}
		if conf.ProxyVia != nil {
			transport.Proxy = http.ProxyURL(conf.ProxyVia)
		}

		timeout := 15 * time.Second

		transport.TLSHandshakeTimeout = timeout
		transport.IdleConnTimeout = timeout
		transport.ExpectContinueTimeout = timeout
		transport.ResponseHeaderTimeout = timeout

		conf.Transport = transport
	}

	return &proxy{conf: conf, e: events}
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
		p.e.Emit(&OnError{Request: evtRequest(r), Err: err})
		http.Error(w, err.Error(), http.StatusServiceUnavailable)
	}

	defer in.Close()

	out, err := net.DialTimeout("tcp", r.Host, 10*time.Second)
	if err != nil {
		p.e.Emit(&OnError{Request: evtRequest(r), Err: err})
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
		p.e.Emit(&OnError{Request: evtRequest(r), Err: err})
	}
}

func (p *proxy) handleHTTP(w http.ResponseWriter, r *http.Request) {
	res, err := p.conf.Transport.RoundTrip(r)
	if err != nil {
		p.e.Emit(&OnError{Request: evtRequest(r), Err: err})
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
		p.e.Emit(&OnError{Request: evtRequest(r), Err: err})
	}
}
