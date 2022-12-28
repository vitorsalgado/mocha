package mocha

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

var _ Reply = (*ProxyReply)(nil)

var (
	_client           = &http.Client{}
	_forbiddenHeaders = []string{
		"Connection",
		"Keep-Alive",
		"ServeHTTP-Authenticate",
		"ServeHTTP-Authorization",
		"TE",
		"Trailer",
		"Transfer-Encoding",
		"Upgrade",
	}
)

const _defaultTimeout = 30 * time.Second

type FromTypes interface{ string | *url.URL | url.URL }

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
}

// From inits a ProxyReply with the given target URL.
func From[T FromTypes](target T) *ProxyReply {
	u := &url.URL{}

	switch e := any(target).(type) {
	case string:
		var err error

		u, err = url.Parse(e)
		if err != nil {
			panic(fmt.Errorf("[reply.From()] unable to url.Parse the value \"%s\". reason=%w", e, err))
		}
	case *url.URL:
		u = e
	case url.URL:
		u = &e
	}

	return &ProxyReply{
		target:               u,
		headers:              make(http.Header),
		proxyHeaders:         make(http.Header),
		proxyHeadersToRemove: make([]string, 0),
		timeout:              _defaultTimeout,
	}
}

// Header sets an extra response header that will be set after proxy target responds.
func (r *ProxyReply) Header(key, value string) *ProxyReply {
	r.headers.Add(key, value)
	return r
}

// Headers sets extra response headers that will be set after proxy target responds.
func (r *ProxyReply) Headers(header http.Header) *ProxyReply {
	r.headers = header.Clone()
	return r
}

// ProxyHeader sets an extra header to be sent to the proxy target.
func (r *ProxyReply) ProxyHeader(key, value string) *ProxyReply {
	r.proxyHeaders.Add(key, value)
	return r
}

// ProxyHeaders sets extra headers to be sent to the proxy target.
func (r *ProxyReply) ProxyHeaders(header http.Header) *ProxyReply {
	r.proxyHeaders = header.Clone()
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

	res, err := _client.Do(req.RawRequest.WithContext(ctx))
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()

	stub := &Stub{Header: make(http.Header, len(res.Header))}

	for _, h := range _forbiddenHeaders {
		res.Header.Del(h)
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

	stub.Body = bytes.NewReader(b)
	stub.Trailer = make(http.Header, len(res.Trailer))

	for key, values := range res.Trailer {
		for _, value := range values {
			stub.Trailer.Add(key, value)
		}
	}

	return stub, nil
}