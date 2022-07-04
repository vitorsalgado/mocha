<h1 id="mocha-top" align="center">Mocha</h1>

<div align="center">
    <a href="#"><img src="docs/logo.png" width="120px" alt="Mocha Logo"></a>
    <p align="center">
        HTTP Mocking Tool for Go
        <br />
    </p>
    <div>
      <a href="https://github.com/vitorsalgado/mocha/actions/workflows/ci.yml">
        <img src="https://github.com/vitorsalgado/mocha/actions/workflows/ci.yml/badge.svg" alt="CI Status" />
      </a>
      <a href="https://codecov.io/gh/vitorsalgado/mocha">
        <img src="https://codecov.io/gh/vitorsalgado/mocha/branch/main/graph/badge.svg?token=XOFUV52P31" alt="Coverage"/>
      </a>
      <a href="#">
        <img alt="GitHub go.mod Go version" src="https://img.shields.io/github/go-mod/go-version/vitorsalgado/mocha">
      </a>
      <a href="https://pkg.go.dev/github.com/vitorsalgado/mocha">
        <img src="https://pkg.go.dev/badge/github.com/vitorsalgado/mocha.svg" alt="Go Reference">
      </a>
    </div>
</div>

## Overview

HTTP server mocking tool for Go.  
**Mocha** creates a real HTTP server and lets you configure response stubs for specific requests based on a set of
criteria. It provides a functional like API that allows you to match any part of a request against a set of matching
functions.

Inspired by [WireMock](https://github.com/wiremock/wiremock) and [Nock](https://github.com/nock/nock).

> Work In Progress

## Installation

```bash
go get -u github.com/vitorsalgado/mocha
```

## Features

- Configure HTTP response stubs for specific requests based on a criteria set.
- Matches request URL, headers, queries, body.
- Stateful matches to create scenarios, mocks for a specific number of calls.
- Response body template.
- Response delays.
- Run in your automated tests.

## Getting Started

```go
package main

import (
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/vitorsalgado/mocha"
	"github.com/vitorsalgado/mocha/matchers"
	"github.com/vitorsalgado/mocha/reply"
)

func Test(t *testing.T) {
	m := mocha.ForTest(t)
	m.Start()

	scoped := m.Mock(mocha.Get(matchers.URLPath("/test")).
		Header("test", matchers.EqualTo("hello")).
		Query("filter", matchers.EqualTo("all")).
		Reply(reply.Created().BodyString("hello world")))

	req, _ := http.NewRequest(http.MethodGet, m.Server.URL+"/test?filter=all", nil)
	req.Header.Add("test", "hello")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}

	body, err := ioutil.ReadAll(res.Body)

	assert.Nil(t, err)
	assert.True(t, scoped.IsDone())
	assert.Equal(t, 201, res.StatusCode)
	assert.Equal(t, string(body), "hello world")
}

```

## Todo

- [ ] Configure mocks with JSON/YAML files
- [ ] CLI
- [ ] Docker
- [ ] Proxy and Record
- [ ] Admin

## Contributing

Check our [Contributing](CONTRIBUTING.md) guide for more details.

## License · [![FOSSA Status](https://app.fossa.com/api/projects/git%2Bgithub.com%2Fvitorsalgado%2Fmocha.svg?type=small)](https://app.fossa.com/projects/git%2Bgithub.com%2Fvitorsalgado%2Fmocha?ref=badge_small)

This project is [MIT Licensed](LICENSE).

<p align="center"><a href="#mocha-top">back to the top</a></p>
