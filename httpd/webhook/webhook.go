package webhook

import (
	"crypto/tls"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"sync"

	"github.com/mitchellh/mapstructure"
	"github.com/rs/zerolog"

	"github.com/vitorsalgado/mocha/v3/httpd"
)

var _ httpd.PostAction = (*WebHook)(nil)

// Name is the WebHook extension identifier used to register it.
const Name = "webhook"

// Input configures an WebHook HTTP request.
// It will be parsed from the parameters map provided during mocking phase.
type Input struct {
	URL            string
	Method         string
	Header         map[string]string
	HeaderTemplate map[string]string
	Body           string
	BodyFile       string `mapstructure:"body_file"`
	SSLVerify      bool   `mapstructure:"ssl_verify"`
	IsTemplate     bool   `mapstructure:"is_template"`
	Transform      Transform
}

// Transform lets users transforms and customize the WebHook Input
// before it is used to build the HTTP request.
// It will be executed everytime WebHook.Run is called.
type Transform func(input *httpd.PostActionInput, args *Input) error

// WebHook is an automated HTTP request sent after a mocked response is served.
type WebHook struct {
	httpClient *http.Client
	rwMutex    sync.RWMutex
}

// New creates a new WebHook instance.
// New should be used to register the WebHook mocha.PostAction within the server instance.
// If you want to use the WebHook in a mock, call Setup.
func New() *WebHook {
	return &WebHook{}
}

// Run calls parses the parameters and calls the configured service.
func (w *WebHook) Run(input *httpd.PostActionInput) error {
	logger := input.App.Logger()
	conf := input.App.Config()

	if input.Args == nil {
		return fmt.Errorf("webhook: arguments must not be nil")
	}

	args := new(Input)
	a, ok := input.Args.(*Input)
	if ok {
		args = a
	} else {
		err := mapstructure.Decode(input.Args, args)
		if err != nil {
			return fmt.Errorf("webhook: error decoding args:\n%w", err)
		}
	}

	if args.Transform != nil {
		err := args.Transform(input, args)
		if err != nil {
			return fmt.Errorf("webhook: failed to Transform args: %w", err)
		}
	}

	if w.httpClient == nil {
		w.rwMutex.Lock()
		defer w.rwMutex.Unlock()

		if w.httpClient == nil {
			if conf.HTTPClientFactory == nil {
				tlsConfig := conf.TLSConfig
				if tlsConfig == nil {
					tlsConfig = &tls.Config{
						InsecureSkipVerify: args.SSLVerify,
						RootCAs:            conf.TLSClientCAs,
						Certificates:       conf.TLSCertificates,
					}
				}

				w.httpClient = &http.Client{Transport: &http.Transport{TLSClientConfig: tlsConfig}}
			} else {
				w.httpClient = conf.HTTPClientFactory()
			}
		}
	}

	method := args.Method
	if len(method) == 0 {
		method = http.MethodGet
	}

	u, err := url.Parse(args.URL)
	if err != nil {
		return fmt.Errorf("webhook: error parsing %s: %w", args.URL, err)
	}

	var body io.Reader
	if len(args.Body) > 0 {
		body = strings.NewReader(args.Body)
	}

	txtURL := u.String()
	logger.Info().
		Str("url", txtURL).
		Msgf("---> WEBHOOK REQUEST RECEIVED %s %s", method, u.Path)

	req, err := http.NewRequestWithContext(input.RawRequest.Context(), method, txtURL, body)
	if err != nil {
		return err
	}

	for k, v := range args.Header {
		req.Header.Add(k, v)
	}

	res, err := w.httpClient.Do(req)
	if err != nil {
		logger.Err(err).
			Str("url", txtURL).
			Msgf("<--- WEBHOOK: %s %s", method, u.Path)

		return fmt.Errorf("webhook: request failed: %w", err)
	}

	level := zerolog.InfoLevel
	if res.StatusCode >= http.StatusBadRequest {
		level = zerolog.WarnLevel
	}

	logger.WithLevel(level).
		Str("url", txtURL).Int("status", res.StatusCode).
		Msgf("<--- WEBHOOK %s %s", method, u.Path)

	return nil
}
