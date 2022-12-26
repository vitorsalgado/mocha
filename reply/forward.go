package reply

import (
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/vitorsalgado/mocha/v3/types"
)

var _ Reply = (*ProxyReply)(nil)

var forbiddenHeaders = []string{
	"Connection",
	"Keep-Alive",
	"ServeHTTP-Authenticate",
	"ServeHTTP-Authorization",
	"TE",
	"Trailers",
	"Transfer-Encoding",
	"Upgrade",
}

type FromTypes interface{ string | *url.URL }

// ProxyReply represents a response stub that will be the response "proxied" from the specified target.
// Use Forward to init a new ProxyReply.
type ProxyReply struct {
	target               *url.URL
	headers              http.Header
	proxyHeaders         http.Header
	proxyHeadersToRemove []string
	trimPrefix           string
	trimSuffix           string
}

// Forward inits a ProxyReply with the given target URL.
func Forward[T FromTypes](target T) *ProxyReply {
	u := &url.URL{}

	switch e := any(target).(type) {
	case string:
		var err error

		u, err = url.Parse(e)
		if err != nil {
			panic(err)
		}

	case *url.URL:
		u = e
	}

	return &ProxyReply{
		target:               u,
		headers:              make(http.Header),
		proxyHeaders:         make(http.Header),
		proxyHeadersToRemove: make([]string, 0),
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

// RemoveProxyHeader removes the given header before sending the request to the proxy target.
func (r *ProxyReply) RemoveProxyHeader(header string) *ProxyReply {
	r.proxyHeadersToRemove = append(r.proxyHeadersToRemove, header)
	return r
}

// RemoveProxyHeaders removes the given headers before sending the request to the proxy target.
func (r *ProxyReply) RemoveProxyHeaders(headers []string) *ProxyReply {
	r.proxyHeadersToRemove = append(r.proxyHeadersToRemove, headers...)
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

func (r *ProxyReply) Prepare() error { return nil }

func (r *ProxyReply) Raw() types.RawValue {
	return types.RawValue{"response_from", map[string]any{
		"target":                  r.target,
		"headers":                 r.headers,
		"proxy_headers":           r.proxyHeaders,
		"proxy_headers_to_remove": r.proxyHeadersToRemove,
		"trim_prefix":             r.trimPrefix,
		"trim_suffix":             r.trimSuffix,
	}}
}

// Build builds a Reply based on the ProxyReply configuration.
func (r *ProxyReply) Build(w http.ResponseWriter, req *types.RequestValues) (*Stub, error) {
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

	res, err := http.DefaultClient.Do(req.RawRequest)
	if err != nil {
		return nil, err
	}

	for _, h := range forbiddenHeaders {
		res.Header.Del(h)
	}

	for key, values := range r.headers {
		for _, value := range values {
			w.Header().Add(key, value)
		}
	}

	for key, values := range res.Header {
		for _, value := range values {
			w.Header().Add(key, value)
		}
	}

	w.WriteHeader(res.StatusCode)

	_, err = io.Copy(w, res.Body)
	if err != nil {
		return nil, err
	}

	return nil, nil
}
