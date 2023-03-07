package mocha

import (
	"context"
	"crypto/tls"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

var _ Reply = (*ProxyReply)(nil)

var _forbiddenHeaders = []string{
	"Connection",
	"Keep-Alive",
	"Proxy-Authenticate",
	"Proxy-Authorization",
	"TE",
	"Trailer",
	"Transfer-Encoding",
	"Upgrade",
}

const _defaultTimeout = 30 * time.Second

type FromTypes interface{ string | *url.URL }

// ProxyReply represents a response stub that will be the response "proxied" from the specified target.
// Use From to init a new ProxyReply.
type ProxyReply struct {
	target               *url.URL
	headers              http.Header
	proxyHeaders         http.Header
	proxyHeadersToRemove []string
	trimPrefix           string
	trimSuffix           string
	timeout              time.Duration
	noFollow             bool
	sslVerify            bool
	httpClient           *http.Client
}

// From creates a ProxyReply with the given target.
func From[T FromTypes](target T) *ProxyReply {
	u := &url.URL{}

	switch e := any(target).(type) {
	case string:
		var err error

		u, err = url.Parse(e)
		if err != nil {
			panic(fmt.Errorf("reply_proxy: unable to parse URL %s. reason=%w", e, err))
		}
	case *url.URL:
		u = e
	}

	return &ProxyReply{
		target:               u,
		headers:              make(http.Header),
		proxyHeaders:         make(http.Header),
		proxyHeadersToRemove: make([]string, 0),
		timeout:              _defaultTimeout,
	}
}

// NoFollow disables following redirects.
func (r *ProxyReply) NoFollow() *ProxyReply {
	r.noFollow = true
	return r
}

// Header sets an extra response header that will be set after the proxy target replies.
func (r *ProxyReply) Header(key, value string) *ProxyReply {
	r.headers.Add(key, value)
	return r
}

// Headers sets extra response headers that will be set after the proxy target replies.
func (r *ProxyReply) Headers(header http.Header) *ProxyReply {
	for k, v := range header {
		for _, vv := range v {
			r.headers.Add(k, vv)
		}
	}

	return r
}

// ForwardHeader sets an extra header to be sent to the proxy target.
func (r *ProxyReply) ForwardHeader(key, value string) *ProxyReply {
	r.proxyHeaders.Add(key, value)
	return r
}

// ProxyHeaders sets extra headers to be sent to the proxy target.
func (r *ProxyReply) ProxyHeaders(header http.Header) *ProxyReply {
	for k, v := range header {
		for _, vv := range v {
			r.proxyHeaders.Add(k, vv)
		}
	}

	return r
}

// RemoveProxyHeaders removes the given header before sending the request to the proxy target.
func (r *ProxyReply) RemoveProxyHeaders(header ...string) *ProxyReply {
	r.proxyHeadersToRemove = append(r.proxyHeadersToRemove, header...)
	return r
}

// TrimPrefix removes the given prefix from the URL before proxying the request.
func (r *ProxyReply) TrimPrefix(prefix string) *ProxyReply {
	r.trimPrefix = prefix
	return r
}

// TrimSuffix removes the given suffix from the URL before proxying the request.
func (r *ProxyReply) TrimSuffix(suffix string) *ProxyReply {
	r.trimSuffix = suffix
	return r
}

// Timeout sets the timeout for the target HTTP request.
// Defaults to 30s.
func (r *ProxyReply) Timeout(timeout time.Duration) *ProxyReply {
	r.timeout = timeout
	return r
}

// SkipSSLVerify skips server certificate verification.
func (r *ProxyReply) SkipSSLVerify() *ProxyReply {
	r.sslVerify = false
	return r
}

// SSLVerify sets if the client should verify server certificate.
func (r *ProxyReply) SSLVerify(v bool) *ProxyReply {
	r.sslVerify = v
	return r
}

// Build builds a Reply based on the ProxyReply configuration.
func (r *ProxyReply) Build(_ http.ResponseWriter, req *RequestValues) (*Stub, error) {
	path := req.RawRequest.URL.Path

	if r.trimPrefix != "" {
		path = strings.TrimPrefix(path, r.trimPrefix)
	}

	if r.trimSuffix != "" {
		path = strings.TrimSuffix(path, r.trimSuffix)
	}

	req.RawRequest.URL.Host = r.target.Host
	req.RawRequest.URL.Scheme = r.target.Scheme
	req.RawRequest.URL.Path = path
	req.RawRequest.Host = r.target.Host
	req.RawRequest.RequestURI = ""

	for _, h := range r.proxyHeadersToRemove {
		req.RawRequest.Header.Del(h)
	}

	for key, values := range r.proxyHeaders {
		for _, value := range values {
			req.RawRequest.Header.Add(key, value)
		}
	}

	ctx, cancel := context.WithTimeout(req.RawRequest.Context(), r.timeout)
	defer cancel()

	res, err := r.httpClient.Transport.RoundTrip(req.RawRequest.WithContext(ctx))
	if err != nil {
		return nil, fmt.Errorf("reply_proxy: error calling target server %s. reason=%w", r.target.String(), err)
	}

	defer res.Body.Close()

	stub := &Stub{Header: make(http.Header)}

	for _, h := range _forbiddenHeaders {
		res.Header.Del(h)
	}

	for k, v := range res.Header {
		for _, vv := range v {
			stub.Header.Add(k, vv)
		}
	}

	for key, values := range r.headers {
		for _, value := range values {
			stub.Header.Add(key, value)
		}
	}

	stub.StatusCode = res.StatusCode

	b, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	stub.Body = b
	stub.Trailer = make(http.Header, len(res.Trailer))

	for key, values := range res.Trailer {
		for _, value := range values {
			stub.Trailer.Add(key, value)
		}
	}

	return stub, nil
}

func (r *ProxyReply) beforeBuild(app *Mocha) error {
	if app.config.HTTPClientFactory == nil {
		r.httpClient = &http.Client{
			Transport: &http.Transport{
				DisableCompression: true,
				TLSClientConfig:    &tls.Config{InsecureSkipVerify: !r.sslVerify},
			}}
	} else {
		h, err := app.config.HTTPClientFactory()
		if err != nil {
			return err
		}

		r.httpClient = h
		if r.httpClient.Transport != nil {
			r.httpClient.Transport.(*http.Transport).TLSClientConfig.InsecureSkipVerify = !r.sslVerify
		}
	}

	if r.noFollow {
		r.httpClient.CheckRedirect = noFollow
	}

	return nil
}

func noFollow(_ *http.Request, _ []*http.Request) error {
	return http.ErrUseLastResponse
}
