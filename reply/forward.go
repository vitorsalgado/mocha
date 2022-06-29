package reply

import (
	"bytes"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/vitorsalgado/mocha/mock"
	"github.com/vitorsalgado/mocha/params"
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

type ProxyReply struct {
	target               string
	headers              http.Header
	proxyHeaders         http.Header
	proxyHeadersToRemove []string
	delay                time.Duration
	trimPrefix           string
	trimSuffix           string
}

func From(target string) *ProxyReply {
	return &ProxyReply{
		target:               target,
		headers:              make(http.Header),
		proxyHeaders:         make(http.Header),
		proxyHeadersToRemove: make([]string, 0),
	}
}

func ProxiedFrom(target string) *ProxyReply {
	return From(target)
}

func (r *ProxyReply) Target(target string) *ProxyReply {
	r.target = target
	return r
}

func (r *ProxyReply) Header(key, value string) *ProxyReply {
	r.headers.Add(key, value)
	return r
}

func (r *ProxyReply) ProxyHeader(key, value string) *ProxyReply {
	r.proxyHeaders.Add(key, value)
	return r
}

func (r *ProxyReply) RemoveProxyHeader(header string) *ProxyReply {
	r.proxyHeadersToRemove = append(r.proxyHeadersToRemove, header)
	return r
}

func (r *ProxyReply) StripPrefix(prefix string) *ProxyReply {
	r.trimPrefix = prefix
	return r
}

func (r *ProxyReply) StripSuffix(suffix string) *ProxyReply {
	r.trimSuffix = suffix
	return r
}

func (r *ProxyReply) Err() error {
	return nil
}

func (r *ProxyReply) Build(req *http.Request, m *mock.Mock, p *params.Params) (*mock.Response, error) {
	t, err := url.Parse(r.target)
	if err != nil {
		return nil, err
	}

	path := req.URL.Path

	if r.trimPrefix != "" {
		path = strings.TrimPrefix(path, r.trimPrefix)
	}

	if r.trimSuffix != "" {
		path = strings.TrimSuffix(path, r.trimSuffix)
	}

	req.URL.Host = t.Host
	req.URL.Scheme = t.Scheme
	req.URL.Path = path

	req.Host = t.Host
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

	response := &mock.Response{
		Status:  res.StatusCode,
		Header:  res.Header,
		Cookies: res.Cookies(),
		Delay:   r.delay,
	}

	buf := &bytes.Buffer{}
	io.Copy(buf, res.Body)
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
