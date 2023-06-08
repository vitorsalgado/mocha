package httpd

import (
	"bytes"
	_ "embed"
	"errors"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"strings"
	"text/template"
	"time"

	"github.com/spf13/viper"
	"github.com/xeipuuv/gojsonschema"

	"github.com/vitorsalgado/mocha/v3/lib"
	"github.com/vitorsalgado/mocha/v3/matcher"
	"github.com/vitorsalgado/mocha/v3/matcher/mbuild"
)

var (
	_ lib.Builder[*HTTPMock, *HTTPMockApp] = (*mockBuilderFromFile)(nil)
	_ lib.Builder[*HTTPMock, *HTTPMockApp] = (*mockBuilderFromBytes)(nil)
	_ error                                = (*jsonSchemaValidationErr)(nil)
)

const (
	_fName                 = "name"
	_fEnabled              = "enabled"
	_fPriority             = "priority"
	_fDelay                = "delay"
	_fRepeat               = "repeat"
	_fExtenders            = "ext"
	_fPostActions          = "post_actions"
	_fPostActionName       = "name"
	_fPostActionParameters = "parameters"

	_fScenario              = "scenario"
	_fScenarioName          = "scenario.name"
	_fScenarioRequiredState = "scenario.required_state"
	_fScenarioNewState      = "scenario.new_state"

	_fRequestScheme   = "request.scheme"
	_fRequestMethod   = "request.method"
	_fRequestURL      = "request.url"
	_fRequestURLMatch = "request.url_match"
	_fRequestURLPath  = "request.path"
	_fRequestQuery    = "request.query"
	_fRequestQueries  = "request.queries"
	_fRequestHeader   = "request.header"
	_fRequestForm     = "request.form_url_encoded"
	_fRequestBody     = "request.body"

	_fResponse         = "response"
	_fResponseStatus   = "response.status"
	_fResponseHeader   = "response.header"
	_fResponseBody     = "response.body"
	_fResponseEncoding = "response.encoding"
	_fResponseBodyFile = "response.body_file"

	_fResponseSequence        = "response_sequence"
	_fResponseSequenceEntries = "response_sequence.responses"
	_fResponseSequenceEnded   = "response_sequence.after_ended"

	_fResponseRandom        = "response_random"
	_fResponseRandomEntries = "response_random.responses"
	_fResponseRandomSeed    = "response_random.seed"

	_fResponseProxy                   = "proxy"
	_fResponseProxyTarget             = "proxy.target"
	_fResponseProxyAdditionalHeader   = "proxy.forward_header"
	_fResponseProxyHeader             = "proxy.header"
	_fResponseProxyHeaderKeysToRemove = "proxy.remove_headers"
	_fResponseProxyTrimPrefix         = "proxy.trim_prefix"
	_fResponseProxyTrimSuffix         = "proxy.trim_suffix"
	_fResponseProxyTimeout            = "proxy.timeout"
	_fResponseProxySSLVerify          = "proxy.ssl_verify"
	_fResponseProxyNoFollow           = "proxy.no_follow"
)

// Response fields
const (
	_resStatus          = "status"
	_resHeader          = "header"
	_resHeaderTemplate  = "header_template"
	_resBody            = "body"
	_resBodyFile        = "body_file"
	_resTemplateEnabled = "template.enabled"
	_resTemplateData    = "template.data"
	_resEncoding        = "encoding"
)

//go:embed schema.json
var _mockJSONSchema []byte

type jsonSchemaValidationErr struct {
	results []gojsonschema.ResultError
}

func (e jsonSchemaValidationErr) Error() string {
	buf := strings.Builder{}
	buf.WriteString("mock definition did not pass the schema validation:\n")
	for _, result := range e.results {
		buf.WriteString(" ")
		buf.WriteString(result.String())
		buf.WriteString("\n")
	}

	return buf.String()
}

type mockBuilderFromFile struct {
	builder  *HTTPMockBuilder
	filename string
}

// FromFile builds a Mock from the given filename.
// It accepts the same extensions from viper.Viper.
// Every mock configuration file is treated as a Go template.
// Check the documentation to see the available data to be used within the template.
// Since every mock file is a Go template by default,
// if you need to define templates for the response URL, header or body, remember to escape it.
// Eg.: body: {{`{{ .Request.Method }}`}}
func FromFile(filename string) lib.Builder[*HTTPMock, *HTTPMockApp] {
	return &mockBuilderFromFile{filename: filename, builder: Request()}
}

func (b *mockBuilderFromFile) Build(app *HTTPMockApp) (mock *HTTPMock, err error) {
	mock, err = b.build(app)
	if err != nil {
		return nil, fmt.Errorf("mock: failed to build mock from file %s.\n%w", b.filename, err)
	}

	return mock, nil
}

func (b *mockBuilderFromFile) build(app *HTTPMockApp) (mock *HTTPMock, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("failed to build mock from file %s.\n%v", b.filename, r)
		}
	}()

	_, err = os.Stat(b.filename)

	if os.IsNotExist(err) {
		return nil, fmt.Errorf("file %s does not exist", b.filename)
	}

	file, err := os.Open(b.filename)
	if err != nil {
		return nil, err
	}

	content, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}

	parts := strings.Split(b.filename, ".")
	ext := parts[len(parts)-1]

	b.builder.SetSource(b.filename)

	return buildMockFromBytes(app, b.builder, content, ext)
}

type mockBuilderFromBytes struct {
	builder *HTTPMockBuilder
	bytes   []byte
	ext     string
}

func FromBytes(b []byte, extension string) lib.Builder[*HTTPMock, *HTTPMockApp] {
	return &mockBuilderFromBytes{bytes: b, builder: Request(), ext: extension}
}

func (b *mockBuilderFromBytes) Build(app *HTTPMockApp) (mock *HTTPMock, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("mock: error building mock from byte array. panic: %v", r)
		}
	}()

	mock, err = buildMockFromBytes(app, b.builder, b.bytes, b.ext)
	if err != nil {
		return nil, fmt.Errorf("mock: error building mock from byte array.\n%w", err)
	}

	return mock, nil
}

// buildReply expects a sub instance of viper.Viper, containing the response definition.
// Fields should not be accessed using "response.", given that the argument v already contains only the
// response.* fields.
func buildReply(vi *viper.Viper) (Reply, error) {
	vi.SetDefault(_resStatus, http.StatusOK)
	vi.SetDefault(_resTemplateEnabled, false)

	res := NewReply()
	res.Status(vi.GetInt(_resStatus))

	for k, v := range vi.GetStringMapString(_resHeader) {
		res.Header(k, v)
	}

	for k, v := range vi.GetStringMapString(_resHeaderTemplate) {
		res.HeaderTemplate(k, v)
	}

	switch vi.GetString(_resEncoding) {
	case "gzip":
		res.Gzip()
	}

	teData := vi.Get(_resTemplateData)
	isTemplateEnabled := vi.GetBool(_resTemplateEnabled) || teData != nil

	if teData != nil {
		res.SetTemplateData(teData)
	}

	if vi.IsSet(_resBodyFile) {
		filename := vi.GetString(_resBodyFile)

		if isTemplateEnabled {
			res.BodyFileWithTemplate(filename)
		} else {
			res.BodyFile(filename)
		}
	} else {
		body := vi.Get(_resBody)
		switch e := body.(type) {
		case string:
			if isTemplateEnabled {
				res.BodyTemplate(e)
			} else {
				res.Body([]byte(e))
			}
		case []map[string]any, map[string]any, bool, float64:
			res.BodyJSON(e)
		}
	}

	return res, nil
}

func buildMockFromBytes(app *HTTPMockApp, builder *HTTPMockBuilder, content []byte, ext string) (*HTTPMock, error) {
	tmpl, err := template.New("").Parse(string(content))
	if err != nil {
		return nil, err
	}

	d := &mockFileData{App: &templateAppWrapper{app}}
	buf := &bytes.Buffer{}
	err = tmpl.Execute(buf, d)
	if err != nil {
		return nil, err
	}

	loader := gojsonschema.NewBytesLoader(_mockJSONSchema)
	validated := false
	if ext == "json" {
		doc := gojsonschema.NewBytesLoader(buf.Bytes())
		result, err := gojsonschema.Validate(loader, doc)
		if err != nil {
			return nil, fmt.Errorf("mock definition did not pass the schema validation: %w", err)
		}

		if len(result.Errors()) > 0 {
			return nil, &jsonSchemaValidationErr{result.Errors()}
		}

		validated = true
	}

	vi := viper.New()
	vi.SetConfigType(ext)

	// Setting defaults before reading the configuration file
	vi.SetDefault(_fEnabled, true)

	err = vi.ReadConfig(buf)
	if err != nil {
		return nil, err
	}

	if !validated {
		doc := gojsonschema.NewGoLoader(vi.AllSettings())
		result, err := gojsonschema.Validate(loader, doc)
		if err != nil {
			return nil, fmt.Errorf("mock definition did not pass the schema validation: %w", err)
		}

		if len(result.Errors()) > 0 {
			return nil, &jsonSchemaValidationErr{result.Errors()}
		}
	}

	err = validateResponseDefinitions(vi)
	if err != nil {
		return nil, err
	}

	// Init building the mock specification
	// --

	// Begin General
	// --

	builder.Name(vi.GetString(_fName))
	builder.Priority(vi.GetInt(_fPriority))
	builder.Enabled(vi.GetBool(_fEnabled))

	if vi.IsSet(_fRepeat) {
		builder.Times(vi.GetInt(_fRepeat))
	}

	if vi.IsSet(_fScenario) {
		builder.ScenarioIs(vi.GetString(_fScenarioName))
		builder.ScenarioStateIs(vi.GetString(_fScenarioRequiredState))
		builder.ScenarioStateWillBe(vi.GetString(_fScenarioNewState))
	}

	if vi.IsSet(_fDelay) {
		var duration time.Duration
		var delay = vi.Get(_fDelay)

		switch v := delay.(type) {
		case string:
			duration, err = time.ParseDuration(v)
			if err != nil {
				return nil, fmt.Errorf("[delay] error parsing delay. %w", err)
			}
		case float64:
			duration = time.Duration(v * float64(time.Millisecond))
		case int:
			duration = time.Duration(float64(v) * float64(time.Millisecond))
		default:
			return nil, fmt.Errorf("[delay] type %T is not supported. supported types: string, number", v)
		}

		builder.Delay(duration)
	}

	if vi.IsSet(_fPostActions) {
		rawPa := vi.Get(_fPostActions)
		postActions, ok := rawPa.([]any)
		if !ok {
			return nil, fmt.Errorf("[post_actions] expected an array of objects. got: %T", rawPa)
		}

		for i, item := range postActions {
			pa, ok := item.(map[string]any)
			if !ok {
				return nil, fmt.Errorf("[post_actions][%d] expected an object. got: %T", i, item)
			}

			builder.PostAction(
				&PostActionDef{Name: pa[_fPostActionName].(string), RawParameters: pa[_fPostActionParameters]})
		}
	}

	// --
	// End General

	// Begin Request
	// --

	if vi.IsSet(_fRequestMethod) {
		mv := vi.Get(_fRequestMethod)

		switch t := mv.(type) {
		case string:
			builder.Method(t)
		default:
			m, err := mbuild.BuildMatcher(mv)
			if err != nil {
				return nil, fmt.Errorf("[request.method] error building matcher %v.\n%w", mv, err)
			}

			builder.MethodMatches(m)
		}
	}

	if vi.IsSet(_fRequestScheme) {
		scheme := vi.Get(_fRequestScheme)

		switch t := scheme.(type) {
		case string:
			builder.Scheme(t)
		default:
			m, err := mbuild.BuildMatcher(scheme)
			if err != nil {
				return nil, fmt.Errorf("[request.scheme] error building matcher %v.\n%w", scheme, err)
			}

			builder.SchemeMatches(m)
		}
	}

	if vi.IsSet(_fRequestURL) {
		uv := vi.Get(_fRequestURL)

		switch e := uv.(type) {
		case string:
			u, err := url.Parse(e)
			if err != nil {
				return nil, fmt.Errorf("[request.url] error parsing url \"%s\".\n%w", e, err)
			}

			if u.IsAbs() {
				builder.URL(matcher.EqualIgnoreCase(e))
			} else {
				builder.URL(matcher.URLPath(u.Path))
			}
		default:
			m, err := mbuild.BuildMatcher(uv)
			if err != nil {
				return nil, fmt.Errorf("[request.url] error building url.\n%w", err)
			}

			builder.URL(m)
		}
	} else if vi.IsSet(_fRequestURLMatch) {
		builder.URL(matcher.Matches(vi.GetString(_fRequestURLMatch)))
	} else if vi.IsSet(_fRequestURLPath) {
		uv := vi.Get(_fRequestURLPath)

		switch e := uv.(type) {
		case string:
			builder.URL(matcher.URLPath(e))
		default:
			m, err := mbuild.BuildMatcher(uv)
			if err != nil {
				return nil, fmt.Errorf("[request.url_path] error building matcher.\n%w", err)
			}

			builder.URLPath(m)
		}
	}

	for k, v := range vi.GetStringMap(_fRequestQuery) {
		m, err := mbuild.TryBuildMatcher(v)
		if err != nil {
			return nil,
				fmt.Errorf("[request.query[%s]] error building matcher.\n%w", k, err)
		}

		builder.Query(k, m)
	}

	for k, v := range vi.GetStringMap(_fRequestQueries) {
		m, err := mbuild.TryBuildMatcher(v)
		if err != nil {
			return nil,
				fmt.Errorf("[request.queries[%s]] error building matcher.\n%w", k, err)
		}

		builder.Queries(k, m)
	}

	for k, v := range vi.GetStringMap(_fRequestHeader) {
		m, err := mbuild.TryBuildMatcher(v)
		if err != nil {
			return nil,
				fmt.Errorf("[request.header[%s]] error building matcher.\n%w", k, err)
		}

		builder.Header(k, m)
	}

	for k, v := range vi.GetStringMap(_fRequestForm) {
		m, err := mbuild.TryBuildMatcher(v)
		if err != nil {
			return nil,
				fmt.Errorf("[request.form[%s]] error building matcher.\n%w", k, err)
		}

		builder.FormField(k, m)
	}

	if vi.IsSet(_fRequestBody) {
		m, err := mbuild.TryBuildMatcher(vi.Get(_fRequestBody))
		if err != nil {
			return nil,
				fmt.Errorf("[request.body] error building matcher.\n%w", err)
		}

		builder.Body(m)
	}

	// --
	// End Request

	// Begin Stub
	// --

	var rep Reply

	if vi.IsSet(_fResponse) {
		rep, err = buildReply(vi.Sub(_fResponse))
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
					fmt.Errorf("[response_random.responses[%d]] error building random response.\n%w", i, err)
			}

			rr, err := buildReply(sub)
			if err != nil {
				return nil,
					fmt.Errorf("[response_random.responses[%d]] building error.\n%w", i, err)
			}

			random.Add(rr)
		}

		rep = random
	} else if vi.IsSet(_fResponseSequence) {
		if !vi.IsSet(_fResponseSequenceEntries) {
			return nil, errors.New("[response_sequence] requires at least one response definition")
		}

		seq := Seq()

		if vi.IsSet(_fResponseSequenceEnded) {
			rr, err := buildReply(vi.Sub(_fResponseSequenceEnded))
			if err != nil {
				return nil,
					fmt.Errorf("[response_response.sequence_ended] building error.\n%w", err)
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

			rr, err := buildReply(sub)
			if err != nil {
				return nil,
					fmt.Errorf("[response_sequence.responses[%d]] building error.\n%w", i, err)
			}

			seq.Add(rr)
		}

		rep = seq
	} else if vi.IsSet(_fResponseProxy) {
		if !vi.IsSet(_fResponseProxyTarget) {
			return nil, errors.New("[proxy.target] is required")
		}

		target := vi.GetString(_fResponseProxyTarget)
		proxy := From(target)

		for k, v := range vi.GetStringMapString(_fResponseProxyAdditionalHeader) {
			proxy.ForwardHeader(k, v)
		}

		for k, v := range vi.GetStringMapString(_fResponseProxyHeader) {
			proxy.Header(k, v)
		}

		for _, k := range vi.GetStringSlice(_fResponseProxyHeaderKeysToRemove) {
			proxy.RemoveProxyHeaders(k)
		}

		if vi.IsSet(_fResponseProxyTrimPrefix) {
			proxy.TrimPrefix(vi.GetString(_fResponseProxyTrimPrefix))
		}

		if vi.IsSet(_fResponseProxyTrimSuffix) {
			proxy.TrimPrefix(vi.GetString(_fResponseProxyTrimSuffix))
		}

		if vi.IsSet(_fResponseProxyTimeout) {
			proxy.Timeout(vi.GetDuration(_fResponseProxyTimeout))
		}

		if vi.IsSet(_fResponseProxySSLVerify) && !vi.GetBool(_fResponseProxySSLVerify) {
			proxy.SkipSSLVerify()
		}

		if vi.IsSet(_fResponseProxyNoFollow) && vi.GetBool(_fResponseProxyNoFollow) {
			proxy.NoFollow()
		}

		rep = proxy
	} else {
		// no response definition found.
		// default to 200 (OK) with nothing more.
		rep = OK()
	}

	builder.Reply(rep)

	// --
	// End Stub

	// User defined handlers

	if len(app.config.MockFileHandlers) > 0 && vi.IsSet(_fExtenders) {
		settings := vi.GetStringMap(_fExtenders)

		for i, handler := range app.config.MockFileHandlers {
			err = handler.Handle(settings, builder)
			if err != nil {
				return nil, fmt.Errorf("custom_field_handler: field handler at index %d failed.\n%w", i, err)
			}
		}
	}

	return builder.Build(app)
}

func validateResponseDefinitions(vi *viper.Viper) error {
	val := 0
	if vi.IsSet(_fResponse) {
		val++
	}
	if vi.IsSet(_fResponseProxy) {
		val++
	}
	if vi.IsSet(_fResponseSequence) {
		val++
	}
	if vi.IsSet(_fResponseRandom) {
		val++
	}

	if val > 1 {
		return fmt.Errorf("more than one response definition set")
	}

	return nil
}
