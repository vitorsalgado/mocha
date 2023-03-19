package webhook

import (
	"bytes"
	"fmt"
	"net/http"
	"net/url"
	"sync"

	"github.com/vitorsalgado/mocha/v3"
	"github.com/vitorsalgado/mocha/v3/internal/preader"
)

var _ mocha.PostAction = (*WebHook)(nil)

const (
	ArgURL          = "url"
	ArgMethod       = "method"
	ArgHeader       = "header"
	ArgBody         = "body"
	ArgBodyTemplate = "template"
)

type Input struct {
	URL    *url.URL
	Method string
	Header http.Header
	Body   []byte
}

type rawInput struct {
	URL    string
	Method string
	Header map[string]string
	Body   string
}

type Transform func(args *Input) error

type WebHook struct {
	Transform  Transform
	HTTPClient *http.Client

	rwMutex sync.RWMutex
}

func (w *WebHook) Run(input *mocha.PostActionInput) error {
	w.rwMutex.Lock()
	defer w.rwMutex.Unlock()

	arg, err := w.buildArgs(input.Args)
	if err != nil {
		return fmt.Errorf("webhook: failed to arguments: %w", err)
	}

	req, err := http.NewRequest(arg.Method, arg.URL.String(), bytes.NewReader(arg.Body))
	if err != nil {
		return err
	}

	if w.Transform != nil {
		err = w.Transform(arg)
		if err != nil {
			return fmt.Errorf("webhook: failed to transform args: %w", err)
		}
	}

	return nil
}

func (w *WebHook) buildArgs(data map[string]any) (*Input, error) {
	args := new(Input)
	pr := preader.New(data)

	rawURL, err := pr.GetStringRequired(ArgURL)
	if err != nil {
		return nil, err
	}

	args.URL, err = url.Parse(rawURL)
	if err != nil {
		return nil, err
	}

	args.Method, err = pr.GetStringRequired(ArgMethod)
	if err != nil {
		return nil, err
	}

	rawMethod, ok := data[ArgMethod]
	if !ok || rawMethod == nil {
		args.Method = http.MethodGet
	} else {
		args.Method, ok = rawMethod.(string)
		if !ok {
			return nil, fmt.Errorf("webhook: method must be a string")
		}
	}

	args.Header = make(http.Header)

	rawHeader, ok := data[ArgHeader]
	if ok && rawHeader != nil {
		h, ok := rawHeader.(map[string]string)
		if !ok {
			return nil, fmt.Errorf("webhook: header must be a string map")
		}

		for k, v := range h {
			args.Header.Add(k, v)
		}
	}

	body, ok := data[ArgBody]
	if !ok || body == nil {
		return args, nil
	}

	return args, nil
}
