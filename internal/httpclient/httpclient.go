package httpclient

import (
	"net"
	"net/http"
	"time"
)

type Options struct {
	Timeout             time.Duration
	TLSHandshakeTimeout time.Duration
	KeepAlive           time.Duration
	DialTimeout         time.Duration
}

func New(options Options) *http.Client {
	return &http.Client{
		Timeout: 5 * time.Second,
		Transport: &http.Transport{
			ResponseHeaderTimeout: options.Timeout,
			TLSHandshakeTimeout:   options.TLSHandshakeTimeout,
			DisableCompression:    true,
			Dial: (&net.Dialer{
				Timeout:   options.DialTimeout,
				KeepAlive: options.KeepAlive,
			}).Dial,
		},
	}
}
