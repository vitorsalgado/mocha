package mocha

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"sync"

	"github.com/vitorsalgado/mocha/v3/matcher"
)

const (
	_extJSON = "json"
)

type (
	ExtSchema struct {
		Name     string
		Enabled  *bool
		Priority int
		Request  *ExtRequest
	}

	ExtRequest struct {
		Method string
		Query  map[string]any
		Header map[string]any
		Body   any
	}

	ExtScenario struct {
		Name          string
		RequiredState string
	}
)

type Loader interface {
	Load(app *Mocha) error
}

var _ Loader = (*FileLoader)(nil)

type FileLoader struct {
}

func (l *FileLoader) Load(app *Mocha) error {
	matches, err := filepath.Glob(app.Config.Pattern)
	if err != nil {
		return err
	}

	wg := sync.WaitGroup{}

	for _, filename := range matches {
		wg.Add(1)

		go func(filename string) {
			defer wg.Done()

			file, err := os.Open(filename)
			if err != nil {
				app.T.Logf("error loading mock file %s. reason=%s", filename, err.Error())
				return
			}

			v := &ExtSchema{}
			err = decode(filename, file, v)
			if err != nil {
				app.T.Logf("error decoding mock file %s. reason=%s", filename, err.Error())
				return
			}

			file.Close()

			m, err := buildExternalMock(app.Config, filename, v)
			if err != nil {
				app.T.Logf("error building mock file %s. reason=%s", filename, err.Error())
				return
			}

			app.AddMocks(m)

			return
		}(filename)
	}

	wg.Wait()

	return nil
}

func decode(filename string, file io.ReadCloser, r *ExtSchema) error {
	parts := strings.Split(filename, "/")
	ext := parts[len(parts)-1]

	switch ext {
	case _extJSON:
		err := json.NewDecoder(file).Decode(r)
		return err
	default:
		err := json.NewDecoder(file).Decode(r)
		return err
	}
}

func buildExternalMock(config *Config, filename string, ext *ExtSchema) (Builder, error) {
	builder := newMockBuilder()

	builder.Name(ext.Name)
	builder.Priority(ext.Priority)
	builder.mock.Source = filename

	if ext.Enabled == nil {
		builder.mock.Enabled = true
	} else {
		builder.mock.Enabled = *ext.Enabled
	}

	if ext.Request != nil {
		builder.Method(ext.Request.Method)

		for k, v := range ext.Request.Query {
			m, err := discoverAndBuildMatcher(v)
			if err != nil {
				return nil, err
			}

			builder.Query(k, m)
		}

		for k, v := range ext.Request.Header {
			m, err := discoverAndBuildMatcher(v)
			if err != nil {
				return nil, err
			}

			builder.Header(k, m)
		}
	}

	if ext.Request.Body != nil {

	}

	return builder, nil
}

func discoverAndBuildMatcher(v any) (m matcher.Matcher, err error) {
	defer func() {
		if r := recover(); r != nil {
			m = nil
			err = fmt.Errorf("panic=%v", r)
			return
		}
	}()

	t := reflect.TypeOf(v)
	switch t.Kind() {
	case reflect.String:
		return matcher.EqualIgnoreCase(v.(string)), nil
	case reflect.Slice:
		val := reflect.ValueOf(v)
		if val.Len() == 0 {
			return nil, fmt.Errorf("array must equal or greather than 1")
		}

		mk, ok := val.Index(0).Interface().(string)
		if !ok {
			return nil, fmt.Errorf("first index must be the matcher name")
		}

		mk = strings.ToLower(mk)

		return buildMatcher(mk, val.Slice(1, val.Len()).Interface())
	default:
		return nil, fmt.Errorf("unsupported type")
	}
}

func extractMultiMatchers(v any) ([]matcher.Matcher, error) {
	a, ok := v.([]any)
	if !ok {
		return nil, fmt.Errorf("only arrays")
	}

	matchers := make([]matcher.Matcher, len(a))

	for _, entry := range a {
		mat, err := discoverAndBuildMatcher(entry)
		if err != nil {
			return nil, err
		}

		matchers = append(matchers, mat)
	}

	return matchers, nil
}

func buildMatcher(key string, args any) (m matcher.Matcher, err error) {
	defer func() {
		if recovery := recover(); recovery != nil {
			m = nil
			err = fmt.Errorf(
				"panic parsing matcher=%s with args=%v. reason=%v",
				key,
				args,
				recovery,
			)

			return
		}
	}()

	switch strings.ToLower(key) {

	case "all", "allof":
		matchers, err := extractMultiMatchers(args)
		if err != nil {
			return nil, err
		}

		return matcher.AllOf(matchers...), nil

	case "any", "anyof":
		matchers, err := extractMultiMatchers(args)
		if err != nil {
			return nil, err
		}

		return matcher.AnyOf(matchers...), nil

	case "contain", "contains":
		return matcher.Contain(args), nil

	case "both":
		matchers, err := extractMultiMatchers(args)
		if err != nil {
			return nil, err
		}

		if len(matchers) != 2 {
			return nil, fmt.Errorf("required 2")
		}

		return matcher.Both(matchers[0], matchers[1]), nil

	case "each":
		m, err := discoverAndBuildMatcher(args)
		if err != nil {
			return nil, err
		}

		return matcher.Each(m), nil

	case "either":
		matchers, err := extractMultiMatchers(args)
		if err != nil {
			return nil, err
		}

		if len(matchers) != 2 {
			return nil, fmt.Errorf("required 2")
		}

		return matcher.Either(matchers[0], matchers[1]), nil

	case "empty":
		return matcher.Empty(), nil

	case "equal", "equals":
		return matcher.Equal(args), nil

	case "equalignorecase", "equalsignorecase", "equalfold":
		str, ok := args.(string)
		if !ok {
			return nil, fmt.Errorf("string required")
		}

		return matcher.EqualIgnoreCase(str), nil

	case "equaljson":
		return matcher.EqualJSON(args), nil

	case "haskey", "havekey":
		str, ok := args.(string)
		if !ok {
			return nil, fmt.Errorf("string required")
		}

		return matcher.HaveKey(str), nil

	case "hasprefix", "startswith":
		str, ok := args.(string)
		if !ok {
			return nil, fmt.Errorf("string required")
		}

		return matcher.HasPrefix(str), nil

	case "hassuffix", "endswith":
		str, ok := args.(string)
		if !ok {
			return nil, fmt.Errorf("string required")
		}

		return matcher.HasSuffix(str), nil

	case "jsonpath":
		a, ok := args.([]any)
		if !ok {
			return nil, fmt.Errorf("array")
		}

		if len(a) != 2 {
			return nil, fmt.Errorf("")
		}

		chain, ok := a[0].(string)
		if !ok {
			return nil, fmt.Errorf("path string")
		}

		m, err := discoverAndBuildMatcher(a[1])
		if err != nil {
			return nil, err
		}

		return matcher.JSONPath(chain, m), nil

	case "len":
		num, ok := args.(float64)
		if !ok {
			return nil, fmt.Errorf("number required")
		}

		return matcher.HaveLen(int(num)), nil

	case "lowercase":
		m, err := discoverAndBuildMatcher(args)
		if err != nil {
			return nil, err
		}

		return matcher.ToLower(m), nil

	case "regex":
		str, ok := args.(string)
		if !ok {
			return nil, fmt.Errorf("string")
		}

		return matcher.Matches(str), nil

	case "not":
		m, err := discoverAndBuildMatcher(args)
		if err != nil {
			return nil, err
		}

		return matcher.Not(m), nil

	case "present":
		return matcher.Present(), nil

	case "split":
		a, ok := args.([]any)
		if !ok {
			return nil, fmt.Errorf("array")
		}

		if len(a) != 2 {
			return nil, fmt.Errorf("")
		}

		separator, ok := a[0].(string)
		if !ok {
			return nil, fmt.Errorf("separator string")
		}

		m, err := discoverAndBuildMatcher(args)
		if err != nil {
			return nil, err
		}

		return matcher.Split(separator, m), nil

	case "trim":
		m, err := discoverAndBuildMatcher(args)
		if err != nil {
			return nil, err
		}

		return matcher.Trim(m), nil

	case "uppercase":
		m, err := discoverAndBuildMatcher(args)
		if err != nil {
			return nil, err
		}

		return matcher.ToUpper(m), nil

	case "urlpath":
		str, ok := args.(string)
		if !ok {
			return nil, fmt.Errorf("")
		}

		return matcher.URLPath(str), nil

	case "xor":
		matchers, err := extractMultiMatchers(args)
		if err != nil {
			return nil, err
		}

		if len(matchers) != 2 {
			return nil, fmt.Errorf("required 2")
		}

		return matcher.XOR(matchers[0], matchers[1]), nil

	default:
		return nil, fmt.Errorf("unknown matcher key=%s", key)
	}
}
