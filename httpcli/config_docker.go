package main

import (
	"os"
	"strings"

	"github.com/vitorsalgado/mocha/v3/httpd"
)

type dockerConfigurer struct {
}

func (c *dockerConfigurer) Apply(conf *httpd.Config) error {
	host := strings.TrimSpace(os.Getenv(_dockerHostEnv))
	if host == "" {
		host = "0.0.0.0:"
	}

	if !strings.HasSuffix(host, ":") {
		host += ":"
	}

	if conf.UseHTTPS {
		conf.Addr = host + "8443"
	} else {
		conf.Addr = host + "8080"
	}

	return nil
}
