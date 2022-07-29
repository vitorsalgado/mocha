<h1 id="mocha-top" align="center">Mocha</h1>

<div align="center">
    <a href="#"><img src="logo.png" width="120px" alt="Mocha Logo"></a>
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
      <a href="https://github.com/vitorsalgado/mocha/releases">
        <img alt="GitHub release (latest by date)" src="https://img.shields.io/github/v/release/vitorsalgado/mocha">
      </a>
      <a href="https://pkg.go.dev/github.com/vitorsalgado/mocha">
        <img src="https://pkg.go.dev/badge/github.com/vitorsalgado/mocha.svg" alt="Go Reference">
      </a>
      <a href="https://goreportcard.com/report/github.com/vitorsalgado/mocha">
        <img src="https://goreportcard.com/badge/github.com/vitorsalgado/mocha" alt="Go Report" />
      </a>
    </div>
</div>

## Overview

HTTP server mocking tool for Go.  
**Mocha** creates a real HTTP server and lets you configure response stubs for HTTP Requests when it matches a set of
matchers.
It provides a functional like API that allows you to match any part of a request against a set of matching
functions that can be composed.

Inspired by [WireMock](https://github.com/wiremock/wiremock) and [Nock](https://github.com/nock/nock).

## Installation

```bash
go get github.com/vitorsalgado/mocha/v2
```

## Features

- Configure HTTP response stubs for specific requests based on a criteria set.
- Matches request URL, headers, queries, body.
- Stateful matches to create scenarios, mocks for a specific number of calls.
- Response body template.
- Response delays.
- Run in your automated tests.

## How It Works

**Mocha** works by creating a real HTTP Server that you can configure response stubs for HTTP requests when they match a
set request matchers. Mock definitions are stored in memory in the server and response will continue to be served as
long as the requests keep passing the configured matchers.  
The basic is workflow for a request is:

- run configured middlewares
- mocha parses the request body based on:
  - custom `RequestBodyParser` configured
  - request content-type
- mock http handler tries to find a mock for the incoming request were all matchers evaluates to true
  - if a mock is found, it will run **post matchers**.
  - if all matchers passes, it will use mock reply implementation to build a response
  - if no mock is found, **it returns an HTTP Status Code 418 (teapot)**.
- after serving a mock response, it will run any `PostAction` configured.

## Getting Started

Usage typically looks like the example below:

```go
func Test_Example(t *testing.T) {
	m := mocha.New(t)
	m.Start()

	scoped := m.AddMocks(mocha.Get(expect.URLPath("/test")).
		Header("test", expect.ToEqual("hello")).
		Query("filter", expect.ToEqual("all")).
		Reply(reply.Created().BodyString("hello world")))

	req, _ := http.NewRequest(http.MethodGet, m.URL()+"/test?filter=all", nil)
	req.Header.Add("test", "hello")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}

	body, err := ioutil.ReadAll(res.Body)

	assert.Nil(t, err)
	assert.True(t, scoped.Called())
	assert.Equal(t, 201, res.StatusCode)
	assert.Equal(t, string(body), "hello world")
}
```

## Configuration

Mocha has two ways to create an instance: `mocha.New()` and `mocha.NewSimple()`.  
`mocha.NewSimple()` creates a new instance with default values for everything.  
`mocha.New(t, ...config)` needs a `mocha.T` implementation and allows to configure the mock server.
You use `testing.T` implementation. Mocha will use this to log useful information for each request match attempt.
Use `mocha.Configure()` or provide a `mocha.Config` to configure the mock server.

## Request Matching

Matchers can be applied to any part of a Request and **Mocha** provides a fluent API to make your life easier.  
See usage examples below:

### Method and URL

```go
m := mocha.New(t)
m.AddMocks(mocha.Request().Method(http.MethodGet).URL(expect.URLPath("/test"))
```

### Header

```go
m := mocha.New(t)
m.AddMocks(mocha.Get(expect.URLPath("/test")).
    Header("test", expect.ToEqual("hello")))
```

### Query

```go
m := mocha.New(t)
m.AddMocks(mocha.Get(expect.URLPath("/test")).
    Query("filter", expect.ToEqual("all")))
```

### Body

**Matching JSON Fields**

```go
m := mocha.New(t)
m.AddMocks(mocha.Post(expect.URLPath("/test")).
    Body(
        expect.JSONPath("name", expect.ToEqual("dev")), expect.JSONPath("ok", expect.ToEqual(true))).
    Reply(reply.OK()))
```

### Form URL Encoded Fields

```go
m.AddMocks(mocha.Post(expect.URLPath("/test")).
    FormField("field1", expect.ToEqual("dev")).
    FormField("field2", expect.ToContain("qa")).
    Reply(reply.OK()))
```

## Replies

You can define a response that should be served once a request is matched.  
**Mocha** provides several ways to configure a reply.  
The built-in reply features are:

- [Basic](#basic-reply)
- [Random Replies](#random-replies)
- [Sequence Replies](#replies-in-sequence)
- [Function](#reply-function)
- [Reply From Proxied Forwarded Request](#reply-from-forwarded-request)
- [Specify Headers](#specifying-headers)
- [Delay Responses](#delay-responses)

Replies are based on the `Reply` interface.  
It's also possible to configure response bodies from templates. **Mocha** uses Go Templates.
Replies usage examples:

### Basic Reply

```go
m := mocha.New(t)
m.AddMocks(mocha.Get(expect.URLPath("/test")).
    Reply(reply.OK())
```

### Replies In Sequence

```go
m := mocha.New(t)
m.AddMocks(mocha.Get(expect.URLPath("/test")).
    Reply(reply.Seq().
	    Add(InternalServerError(), BadRequest(), OK(), NotFound())))
```

### Random Replies

```go
m := mocha.New(t)
m.AddMocks(mocha.Get(expect.URLPath("/test")).
    Reply(reply.Rand().
		Add(reply.BadRequest(), reply.OK(), reply.Created(), reply.InternalServerError())))
```

### Reply Function

```go
m := mocha.New(t)
m.AddMocks(mocha.Get(expect.URLPath("/test")).
    ReplyFunction(func(r *http.Request, m mocha.M, p params.P) (*mocha.Response, error) {
        return &mocha.Response{Status: http.StatusAccepted}, nil
    }))
```

### Reply From Forwarded Request

**reply.From** will forward the request to the given destination and serve the response from the forwarded server.  
It`s possible to add extra headers to the request and the response and also remove unwanted headers.

```go
m := mocha.New(t)
m.AddMocks(mocha.Get(expect.URLPath("/test")).
    Reply(reply.From("http://example.org").
		ProxyHeader("x-proxy", "proxied").
		RemoveProxyHeader("x-to-be-removed").
		Header("x-res", "response"))
```

### Body Template

**Mocha** comes with a built-in template parser based on Go Templates.  
To serve a response body from a template, follow the example below:

```go
templateFile, _ := os.Open("template.tmpl"))
content, _ := ioutil.ReadAll(templateFile)

m := mocha.New(t)
m.AddMocks(mocha.Get(expect.URLPath("/test")).
    Reply(reply.
        OK().
        BodyTemplate(reply.NewTextTemplate().
            FuncMap(template.FuncMap{"trim": strings.TrimSpace}).
            Template(string(content))).
        Model(data))
```

### Specifying Headers

```go
m := mocha.New(t)
m.AddMocks(mocha.Get(expect.URLPath("/test")).
    Reply(reply.OK().Header("test", "test-value"))
```

### Delay Responses

You can configure a delay to responses to simulate timeouts, slow requests and any other timing related scenarios.  
See the example below:

```go
delay := time.Duration(1250) * time.Millisecond

m.AddMocks(Get(expect.URLPath("/test")).
    Reply(reply.
        OK().
        Delay(delay)))
```

## Assertions

### Mocha Instance

Mocha instance provides methods to assert if associated mocks were called or not, how many times they were called,
allows you to enable/disable then and so on.  
The available assertion methods on mocha instance are:

- AssertCalled: asserts that all associated mocks were called at least once.
- AssertNotCalled: asserts that associated mocks were **not** called.
- AssertHits: asserts that the sum of calls is equal to the expected value.

### Scope

Mocha instance method `AddMocks` returns a `Scoped` instance that holds all mocks created.  
`Scopes` allows you control related mocks, enabling/disabling, checking if they were called or not. Scoped instance also
provides **assertions** to facility **tests** verification.
See below the available test assertions:

- AssertCalled: asserts that all associated mocks were called at least once.
- AssertNotCalled: asserts that associated mocks were **not** called.

## Matchers

Mocha provides several matcher functions to facilitate request matching and verification.
See the package `expect` for more details.  
You can create custom matchers using these two approaches:

- create a `expect.Matcher` struct
- use the function `expect.Func` providing a function with the following
  signature: `func(v any, a expect.Args) (bool, error)`

### Matcher Composition

It's possible to compose multiple matchers.  
Every matcher has a `.And()`, `.Or()` and a `.Xor()` that allows composing multiple matchers.  
See the example below:

```go
expect.ToEqual("test").And(expect.ToContain("t"))
```

### BuiltIn Matchers

| Matcher      | Description                                                                                              |
| ------------ | -------------------------------------------------------------------------------------------------------- |
| AllOf        | Returns true when all given matchers returns true                                                        |
| AnyOf        | Returns true when any given matchers returns true                                                        |
| Both         | Returns true when both matchers returns true                                                             |
| ToContain    | Returns true when expected value is contained on the request value                                       |
| Either       | Returns true when any matcher returns true                                                               |
| ToBeEmpty    | Returns true when request value is empty                                                                 |
| ToEqual      | Returns true when values are equal                                                                       |
| ToEqualFold  | Returns true when string values are equal, ignoring case                                                 |
| ToEqualJSON  | Returns true when the expected struct represents a JSON value                                            |
| Func         | Wraps a function to create a inline matcher                                                              |
| ToHaveKey    | Returns true if the JSON key in the given path is present                                                |
| ToHavePrefix | Returns true if the matcher argument starts with the given prefix                                        |
| ToHaveSuffix | Returns true when matcher argument ends with the given suffix                                            |
| JSONPath     | Applies the provided matcher to the JSON field value in the given path                                   |
| ToHaveLen    | Returns true when matcher argument length is equal to the expected value                                 |
| LowerCase    | Lower case matcher string argument before submitting it to provided matcher.                             |
| UpperCase    | Upper case matcher string argument before submitting it to provided matcher                              |
| ToMatchExpr  | Returns true then the given regular expression matches matcher argument                                  |
| Not          | Negates the provided matcher                                                                             |
| Peek         | Will return the result of the given matcher, after executing the provided function                       |
| ToBePresent  | Checks if matcher argument contains a value that is not nil or the zero value for the argument type      |
| Repeat       | Returns true if total request hits for current mock is equal or lower total the provided max call times. |
| Trim         | Trims' spaces of matcher argument before submitting it to the given matcher                              |
| URLPath      | Returns true if request URL path is equal to the expected path, ignoring case                            |
| XOR          | Exclusive "or" matcher                                                                                   |

---

## Future Plans

- [ ] Configure mocks with JSON/YAML files
- [ ] CLI
- [ ] Docker
- [ ] Proxy and Record

## Contributing

Check our [Contributing](CONTRIBUTING.md) guide for more details.

## License

[![FOSSA Status](https://app.fossa.com/api/projects/git%2Bgithub.com%2Fvitorsalgado%2Fmocha.svg?type=shield)](https://app.fossa.com/projects/git%2Bgithub.com%2Fvitorsalgado%2Fmocha?ref=badge_shield)

This project is [MIT Licensed](LICENSE).

<p align="center"><a href="#mocha-top">back to the top</a></p>
