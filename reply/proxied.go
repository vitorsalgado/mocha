package reply

import (
	"bytes"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/vitorsalgado/mocha/v3/params"
)

var forbiddenHeaders = []string{
	"Connection",
	"Keep-Alive",
	"Proxy-Authenticate",
	"Proxy-Authorization",
	"TE",
	"Trailers",
	"Transfer-Encoding",
	"Upgrade",
}

type FromTypes interface{ string | *url.URL }

// ProxyReply represents a response stub that will be the response "proxied" from the specified target.
// Use From or ProxyFrom to init a new ProxyReply.
type ProxyReply struct {
	target               *url.URL
	headers              http.Header
	proxyHeaders         http.Header
	proxyHeadersToRemove []string
	delay                time.Duration
	trimPrefix           string
	trimSuffix           string
}

// From inits a ProxyReply with the given target. It will parse the raw URL target to an url.URL.
// It panics if the given target cannot be parsed to a valid url.URL.
func From(target string) *ProxyReply {
	u, err := url.Parse(target)
	if err != nil {
		panic(err)
	}

	return ProxiedFrom(u)
}

// ProxiedFrom inits a ProxyReply with the given target URL.
func ProxiedFrom(target *url.URL) *ProxyReply {
	return &ProxyReply{
		target:               target,
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

// ProxyHeader sets an extra header to be sent to the proxy target.
func (r *ProxyReply) ProxyHeader(key, value string) *ProxyReply {
	r.proxyHeaders.Add(key, value)
	return r
}

// RemoveProxyHeader removes the given before sending the request to the proxy target.
func (r *ProxyReply) RemoveProxyHeader(header string) *ProxyReply {
	r.proxyHeadersToRemove = append(r.proxyHeadersToRemove, header)
	return r
}

// StripPrefix removes the given prefix from the URL before proxying the request.
func (r *ProxyReply) StripPrefix(prefix string) *ProxyReply {
	r.trimPrefix = prefix
	return r
}

// StripSuffix removes the given suffix from the URL before proxying the request.
func (r *ProxyReply) StripSuffix(suffix string) *ProxyReply {
	r.trimSuffix = suffix
	return r
}

// Build builds a Reply based on the ProxyReply configuration.
func (r *ProxyReply) Build(req *http.Request, _ M, _ params.P) (*Response, error) {
	path := req.URL.Path

	if r.trimPrefix != "" {
		path = strings.TrimPrefix(path, r.trimPrefix)
	}

	if r.trimSuffix != "" {
		path = strings.TrimSuffix(path, r.trimSuffix)
	}

	req.URL.Host = r.target.Host
	req.URL.Scheme = r.target.Scheme
	req.URL.Path = path

	req.Host = r.target.Host
	req.RequestURI = ""

	for _, h := range r.proxyHeadersToRemove {
		req.Header.Del(h)
	}

	for key, values := range r.proxyHeaders {
		for _, value := range values {
			req.Header.Add(key, value)
		}
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	response := &Response{
		Status:  res.StatusCode,
		Header:  res.Header,
		Cookies: res.Cookies(),
		Delay:   r.delay,
	}

	buf := &bytes.Buffer{}
	_, err = io.Copy(buf, res.Body)
	if err != nil {
		return nil, err
	}

	response.Body = buf

	for _, h := range forbiddenHeaders {
		response.Header.Del(h)
	}

	for key, values := range r.headers {
		for _, value := range values {
			response.Header.Add(key, value)
		}
	}

	return response, nil
}
