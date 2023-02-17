package mocha

import (
	"errors"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"path"
	"strings"
	"time"

	"github.com/spf13/viper"

	"github.com/vitorsalgado/mocha/v3/matcher"
	"github.com/vitorsalgado/mocha/v3/matcher/mbuild"
)

// Mock configuration fields
const (
	_fName                       = "name"
	_fEnabled                    = "enabled"
	_fPriority                   = "priority"
	_fDelayInMs                  = "delay_ms"
	_fScenarioName               = "scenario.name"
	_fScenarioRequiredState      = "scenario.required_state"
	_fScenarioNewState           = "scenario.new_state"
	_fRequestMethod              = "request.method"
	_fRequestURL                 = "request.url"
	_fRequestURLMatch            = "request.url_match"
	_fRequestURLPath             = "request.url_path"
	_fRequestURLPathMatch        = "request.url_path_match"
	_fRequestQuery               = "request.query"
	_fRequestHeader              = "request.header"
	_fRequestForm                = "request.form"
	_fRequestBody                = "request.body"
	_fResponse                   = "response"
	_fResponseStatus             = "response.status"
	_fResponseHeader             = "response.header"
	_fResponseBody               = "response.body"
	_fResponseEncoding           = "response.encoding"
	_fResponseBodyFile           = "response.body_file"
	_fResponseTemplateEnabled    = "response.template_enabled"
	_fResponseTemplateModel      = "response.template_model"
	_fResponseSequence           = "response_sequence"
	_fResponseSequenceEntries    = "response_sequence.responses"
	_fResponseSequenceAfterEnded = "response_sequence.after_ended"
	_fResponseRandom             = "response_random"
	_fResponseRandomEntries      = "response_random.responses"
	_fResponseRandomSeed         = "response_random.seed"
)

// Response fields
const (
	_resStatus          = "status"
	_resHeader          = "header"
	_resBody            = "body"
	_resBodyFile        = "body_file"
	_resTemplateEnabled = "template_enabled"
	_resTemplateModel   = "template_model"
	_resEncoding        = "encoding"
)

type mockExternalBuilder struct {
	filename string
	builder  *MockBuilder
}

func MockFromFile(filename string) Builder {
	return &mockExternalBuilder{filename: filename, builder: Request()}
}

func (b *mockExternalBuilder) Build(app *Mocha) (mock *Mock, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("[panic] building external mock. reason: %v", r)
		}
	}()

	vi := viper.New()
	vi.SetConfigFile(b.filename)
	vi.SetDefault(_fEnabled, true)

	err = vi.ReadInConfig()
	if err != nil {
		return nil, err
	}

	// Init building the mock specification
	// --

	// Begin General
	// --

	b.builder.mock.Source = b.filename

	b.builder.Name(vi.GetString(_fName))
	b.builder.Priority(vi.GetInt(_fPriority))
	b.builder.Enable(vi.GetBool(_fEnabled))
	b.builder.ScenarioIs(vi.GetString(_fScenarioName))
	b.builder.ScenarioStateIs(vi.GetString(_fScenarioRequiredState))
	b.builder.ScenarioStateWillBe(vi.GetString(_fScenarioNewState))

	if vi.IsSet(_fDelayInMs) {
		b.builder.Delay(time.Duration(vi.GetInt64(_fDelayInMs)) * time.Millisecond)
	}

	// --
	// End General

	// Begin Request
	// --

	if vi.IsSet(_fRequestMethod) {
		mv := vi.Get(_fRequestMethod)

		switch t := mv.(type) {
		case string:
			b.builder.Method(t)
		case []string:
			b.builder.MethodMatches(matcher.Some(t))
		default:
			m, err := mbuild.BuildMatcher(mv)
			if err != nil {
				return nil, fmt.Errorf("[request.method] error building matcher %v. %w", mv, err)
			}

			b.builder.MethodMatches(m)
		}
	}

	if vi.IsSet(_fRequestURL) {
		uv := vi.Get(_fRequestURL)
		urlConv, ok := uv.(string)

		if ok {
			u, err := url.Parse(urlConv)
			if err != nil {
				return nil, fmt.Errorf("[request.url] error parsing url \"%s\". %w", urlConv, err)
			}

			if u.IsAbs() {
				b.builder.URL(matcher.EqualIgnoreCase(urlConv))
			} else {
				b.builder.URL(matcher.URLPath(u.Path))
			}
		} else {
			m, err := mbuild.BuildMatcher(uv)
			if err != nil {
				return nil, fmt.Errorf("[request.url] error building url. %w", err)
			}

			b.builder.URL(m)
		}
	} else if vi.IsSet(_fRequestURLMatch) {
		b.builder.URL(matcher.Matches(vi.GetString(_fRequestURLMatch)))
	} else if vi.IsSet(_fRequestURLPath) {
		uv := vi.Get(_fRequestURLPath)
		urlConv, ok := uv.(string)
		if ok {
			b.builder.URL(matcher.URLPath(urlConv))
		} else {
			m, err := mbuild.BuildMatcher(uv)
			if err != nil {
				return nil, fmt.Errorf("[request.url_path] error building matcher. %w", err)
			}

			b.builder.URLPath(m)
		}
	} else if vi.IsSet(_fRequestURLPathMatch) {
		b.builder.URLPath(matcher.Matches(vi.GetString(_fRequestURLPathMatch)))
	}

	for k, v := range vi.GetStringMap(_fRequestQuery) {
		m, err := mbuild.BuildMatcher(v)
		if err != nil {
			return nil,
				fmt.Errorf("[request.query[%s]] error building matcher. %w", k, err)
		}

		b.builder.Query(k, m)
	}

	for k, v := range vi.GetStringMap(_fRequestHeader) {
		m, err := mbuild.BuildMatcher(v)
		if err != nil {
			return nil,
				fmt.Errorf("[request.header[%s]] error building matcher. %w", k, err)
		}

		b.builder.Header(k, m)
	}

	for k, v := range vi.GetStringMap(_fRequestForm) {
		m, err := mbuild.BuildMatcher(v)
		if err != nil {
			return nil,
				fmt.Errorf("[request.form[%s]] error building matcher. %w", k, err)
		}

		b.builder.FormField(k, m)
	}

	if vi.IsSet(_fRequestBody) {
		m, err := mbuild.BuildMatcher(vi.Get(_fRequestBody))
		if err != nil {
			return nil,
				fmt.Errorf("[request.body] error building matcher. %w", err)
		}

		b.builder.Body(m)
	}

	// --
	// End Request

	// Begin Stub
	// --

	var rep Reply

	if vi.IsSet(_fResponse) {
		rep, err = b.buildReply(vi.Sub(_fResponse))
		if err != nil {
			return nil, err
		}
	} else if vi.IsSet(_fResponseRandom) {
		if !vi.IsSet(_fResponseRandomEntries) {
			return nil, errors.New("[response_random.responses] requires at least one response definition")
		}

		var random *RandomReply

		if vi.IsSet(_fResponseRandomSeed) {
			random = RandWith(rand.New(rand.NewSource(vi.GetInt64(_fResponseRandomSeed))))
		} else {
			random = Rand()
		}

		entries, ok := vi.Get(_fResponseRandomEntries).([]any)
		if !ok {
			return nil, errors.New("[response_random.response] requires an array of responses")
		}

		for i, r := range entries {
			sub := viper.New()
			err = sub.MergeConfigMap(r.(map[string]any))
			if err != nil {
				return nil,
					fmt.Errorf("[response_random.responses[%d]] error building random response. %w", i, err)
			}

			rr, err := b.buildReply(sub)
			if err != nil {
				return nil,
					fmt.Errorf("[response_random.responses[%d]] building error. %w", i, err)
			}

			random.Add(rr)
		}

		rep = random
	} else if vi.IsSet(_fResponseSequence) {
		if !vi.IsSet(_fResponseSequenceEntries) {
			return nil, errors.New("[response_sequence] requires at least one response definition")
		}

		seq := Seq()

		if vi.IsSet(_fResponseSequenceAfterEnded) {
			rr, err := b.buildReply(vi.Sub(_fResponseSequenceAfterEnded))
			if err != nil {
				return nil,
					fmt.Errorf("[response_response.after_ended] building error. %w", err)
			}

			seq.OnSequenceEnded(rr)
		}

		entries, ok := vi.Get(_fResponseSequenceEntries).([]any)
		if !ok {
			return nil, errors.New("[response_sequence.response] requires an array of responses")
		}

		for i, r := range entries {
			sub := viper.New()
			err = sub.MergeConfigMap(r.(map[string]any))
			if err != nil {
				return nil, err
			}

			rr, err := b.buildReply(sub)
			if err != nil {
				return nil,
					fmt.Errorf("[response_sequence.responses[%d]] building error. %w", i, err)
			}

			seq.Add(rr)
		}

		rep = seq
	} else {
		// no response definition found.
		// default to 200 (OK) with nothing more.
		rep = OK()
	}

	b.builder.Reply(rep)

	// --
	// End Stub

	// User defined handlers

	if len(app.config.MockFileHandlers) > 0 {
		settings := vi.AllSettings()

		for i, handler := range app.config.MockFileHandlers {
			err = handler.Handle(settings, b.builder)
			if err != nil {
				return nil, fmt.Errorf("custom field parser [%d] failed. %w", i, err)
			}
		}
	}

	return b.builder.Build(app)
}

// buildReply expects a sub instance of viper.Viper, containing the response definition.
// Fields should not be accessed using "response.", given that the argument v already contains only the
// response.* fields.
func (b *mockExternalBuilder) buildReply(v *viper.Viper) (Reply, error) {
	v.SetDefault(_resStatus, http.StatusOK)
	v.SetDefault(_resTemplateEnabled, false)

	res := NewReply()
	res.Status(v.GetInt(_resStatus))

	for k, v := range v.GetStringMapString(_resHeader) {
		res.Header(k, v)
	}

	switch v.GetString(_resEncoding) {
	case "gzip":
		res.Gzip()
	}

	if v.IsSet(_resBodyFile) {
		filename := v.GetString(_resBodyFile)
		if !path.IsAbs(filename) {
			dirs := strings.Split(b.filename, "/")
			filename = path.Join(strings.Join(dirs[:len(dirs)-1], "/"), filename)
		}

		file, err := os.Open(filename)
		if err != nil {
			return nil, fmt.Errorf(
				"[%s] error opening file. %w",
				filename,
				err,
			)
		}

		defer file.Close()

		b, err := io.ReadAll(file)
		if err != nil {
			return nil, fmt.Errorf(
				"[%s] error reading file content. %w",
				filename,
				err,
			)
		}

		if v.GetBool(_resTemplateEnabled) {
			res.BodyTemplate(string(b), v.Get(_resTemplateEnabled))
		} else {
			res.Body(b)
		}
	} else {
		b := v.Get(_resBody)
		switch e := b.(type) {
		case string:
			if v.GetBool(_resTemplateEnabled) {
				res.BodyTemplate(e, v.Get(_resTemplateModel))
			} else {
				res.Body([]byte(e))
			}
		case []map[string]any, map[string]any, bool, float64:
			res.BodyJSON(e)
		}
	}

	return res, nil
}
