package mocha

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"reflect"
	"regexp"
	"strings"
	"time"

	"github.com/vitorsalgado/mocha/v3/matcher"
	"github.com/vitorsalgado/mocha/v3/mod"
	"github.com/vitorsalgado/mocha/v3/reply"
)

// Loader is the interface that defines custom Mock loaders.
// Usually, it is used to load external mocks, like from the file system.
type Loader interface {
	Load(app *Mocha) error
}

var _RegExpAbsolutePath = regexp.MustCompile("^[a-zA-Z][a-zA-Z\\d+\\-.]*?:")

func buildExternalMock(source string, ext *mod.ExtMock) (b Builder, err error) {
	defer func() {
		if recovery := recover(); recovery != nil {
			b = nil
			err = fmt.Errorf("panic %v", err)
		}
	}()

	// Init building the mock specification
	// --

	builder := Request()

	// Begin General
	// --

	builder.Name(ext.Name)
	builder.Priority(ext.Priority)

	builder.mock.Source = source

	if ext.Enabled == nil {
		builder.mock.Enabled = true
	} else {
		builder.mock.Enabled = *ext.Enabled
	}

	builder.
		ScenarioIs(ext.Scenario.Name).
		ScenarioStateIs(ext.Scenario.RequiredState).
		ScenarioStateWillBe(ext.Scenario.NewState)

	builder.Delay(time.Duration(ext.DelayInMs) * time.Millisecond)

	// --
	// End General

	// Begin Request
	// --

	if ext.Request.Method != "" {
		builder.Method(ext.Request.Method)
	}

	if ext.Request.URL != nil {
		urlConv, ok := ext.Request.URL.(string)
		if ok {
			if _RegExpAbsolutePath.MatchString(urlConv) {
				builder.URL(matcher.EqualIgnoreCase(urlConv))
			} else {
				q := strings.Index(urlConv, "?")
				if q < 0 {
					builder.URL(matcher.URLPath(urlConv))
				} else {
					builder.URL(matcher.URLPath(urlConv[0:q]))
				}
			}
		} else {
			m, err := matcher.BuildMatcher(ext.Request.URL)
			if err != nil {
				return nil, fmt.Errorf("[request.url] matchers error. %w", err)
			}

			builder.URL(m)
		}
	} else if ext.Request.URLMatch != "" {
		builder.URL(matcher.Matches(ext.Request.URLMatch))
	} else if ext.Request.URLPath != nil {
		urlConv, ok := ext.Request.URL.(string)
		if ok {
			builder.URL(matcher.URLPath(urlConv))
		} else {
			m, err := matcher.BuildMatcher(ext.Request.URL)
			if err != nil {
				return nil, err
			}

			builder.URLPath(m)
		}
	} else if ext.Request.URLPathMatch != "" {
		builder.URLPath(matcher.Matches(ext.Request.URLPathMatch))
	}

	for k, v := range ext.Request.Query {
		m, err := matcher.BuildMatcher(v)
		if err != nil {
			return nil,
				fmt.Errorf("[request.query[%s]] matchers error. %w", k, err)
		}

		builder.Query(k, m)
	}

	for k, v := range ext.Request.Header {
		m, err := matcher.BuildMatcher(v)
		if err != nil {
			return nil,
				fmt.Errorf("[request.header[%s]] matchers error. %w", k, err)
		}

		builder.Header(k, m)
	}

	if ext.Request.Body != nil {
		m, err := matcher.BuildMatcher(ext.Request.Body)
		if err != nil {
			return nil,
				fmt.Errorf("[request.body] matchers error. %w", err)
		}

		builder.Body(m)
	}

	// --
	// End Request

	// Begin Response
	// --

	var rep reply.Reply

	if ext.Response != nil {
		rep, err = buildResponse(ext, ext.Response)
		if err != nil {
			return nil, err
		}
	} else if ext.RandomResponse != nil {
		random := reply.Rand()
		for i, r := range ext.RandomResponse.Responses {
			rr, err := buildResponse(ext, &r)
			if err != nil {
				return nil,
					fmt.Errorf("[random_response.responses[%d]] building error. %w", i, err)
			}

			random.Add(rr)
		}

		rep = random
	} else if ext.SequenceResponse != nil {
		seq := reply.Seq()

		if ext.SequenceResponse.AfterEnded != nil {
			rr, err := buildResponse(ext, ext.SequenceResponse.AfterEnded)
			if err != nil {
				return nil,
					fmt.Errorf("[sequence_response.after_ended] building error. %w", err)
			}

			seq.AfterEnded(rr)
		}

		for i, r := range ext.SequenceResponse.Responses {
			rr, err := buildResponse(ext, &r)
			if err != nil {
				return nil,
					fmt.Errorf("[sequence_response.responses[%d]] building error. %w", i, err)
			}

			seq.Add(rr)
		}

		rep = seq
	} else {
		// no response definition found.
		// default to 200 (OK) with nothing more.
		rep = reply.OK()
	}

	builder.Reply(rep)

	// --
	// End Response

	return builder, nil
}

func buildResponse(ext *mod.ExtMock, response *mod.ExtMockResponse) (reply.Reply, error) {
	res := reply.New()
	res.Status(valueOr(response.Status, http.StatusOK))

	for k, v := range response.Header {
		res.Header(k, v)
	}

	if response.BodyFile != "" {
		file, err := os.Open(response.BodyFile)
		if err != nil {
			return nil, fmt.Errorf(
				"[%s] error opening file. %w",
				response.BodyFile,
				err,
			)
		}

		defer file.Close()

		b, err := io.ReadAll(file)
		if err != nil {
			return nil, fmt.Errorf(
				"[%s] error reading file content. %w",
				response.BodyFile,
				err,
			)
		}

		if response.Template {
			res.BodyTemplate(string(b)).TemplateModel(response.TemplateModel)
		} else {
			res.Body(b)
		}
	} else {
		switch e := response.Body.(type) {
		case string:
			if response.Template {
				res.BodyTemplate(e).TemplateModel(response.TemplateModel)
			} else {
				res.Body([]byte(e))
			}
			break
		case nil:
			break
		default:
			res.BodyJSON(ext.Request.Body)
		}
	}

	return res, nil
}

func valueOr[V any](v V, d V) V {
	if any(v) == reflect.Zero(reflect.TypeOf(v)).Interface() {
		return d
	}

	return v
}
