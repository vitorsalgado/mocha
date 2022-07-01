package httpclient

import (
	"net/http"
	"time"
)

// Options represents available internal http.Client options
type Options struct {
	Timeout             time.Duration
	TLSHandshakeTimeout time.Duration
	KeepAlive           time.Duration
	DialTimeout         time.Duration
}

// New creates a new configurable http.Client instance
func New(options Options) *http.Client {
	return &http.Client{
		Timeout: 5 * time.Second,
		Transport: &http.Transport{
			ResponseHeaderTimeout: options.Timeout,
			TLSHandshakeTimeout:   options.TLSHandshakeTimeout,
			DisableCompression:    true,
		},
	}
}
