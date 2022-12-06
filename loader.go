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
	"github.com/vitorsalgado/mocha/v3/matcher/asm"
	"github.com/vitorsalgado/mocha/v3/mod"
	"github.com/vitorsalgado/mocha/v3/reply"
)

type Loader interface {
	Load(app *Mocha) error
}

var _RegExpAbsolutePath = regexp.MustCompile("^[a-zA-Z][a-zA-Z\\d+\\-.]*?:")

func buildExternalMock(source string, ext *mod.ExternalSchema) (b Builder, err error) {
	defer func() {
		if recovery := recover(); recovery != nil {
			b = nil
			err = fmt.Errorf("%v", err)
		}
	}()

	// build request

	builder := Request()

	builder.Name(ext.Name)
	builder.Priority(ext.Priority)

	builder.mock.Source = source

	if ext.Enabled == nil {
		builder.mock.Enabled = true
	} else {
		builder.mock.Enabled = *ext.Enabled
	}

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
			m, err := asm.BuildMatcher(ext.Request.URL)
			if err != nil {
				return nil, err
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
			m, err := asm.BuildMatcher(ext.Request.URL)
			if err != nil {
				return nil, err
			}

			builder.URLPath(m)
		}
	} else if ext.Request.URLPathMatch != "" {
		builder.URLPath(matcher.Matches(ext.Request.URLPathMatch))
	}

	for k, v := range ext.Request.Query {
		m, err := asm.BuildMatcher(v)
		if err != nil {
			return nil, err
		}

		builder.Query(k, m)
	}

	for k, v := range ext.Request.Header {
		m, err := asm.BuildMatcher(v)
		if err != nil {
			return nil, err
		}

		builder.Header(k, m)
	}

	if ext.Request.Body != nil {
		m, err := asm.BuildMatcher(ext.Request.Body)
		if err != nil {
			return nil, err
		}

		builder.Body(m)
	}

	builder.
		ScenarioIs(ext.Scenario.Name).
		ScenarioStateIs(ext.Scenario.RequiredState).
		ScenarioStateWillBe(ext.Scenario.NewState)

	builder.Delay(time.Duration(ext.DelayInMs) * time.Millisecond)

	if ext.Repeat != nil {
		builder.Repeat(*ext.Repeat)
	}

	// build response

	var rep reply.Reply

	if ext.Response != nil {
		rep, err = buildResponse(ext, ext.Response)
		if err != nil {
			return nil, err
		}
	} else if ext.RandomResponse != nil {
		random := reply.Rand()
		for _, r := range ext.RandomResponse.Responses {
			rr, err := buildResponse(ext, &r)
			if err != nil {
				return nil, err
			}

			random.Add(rr)
		}

		rep = random
	} else if ext.SequenceResponse != nil {
		seq := reply.Seq()

		if ext.SequenceResponse.AfterEnded != nil {
			rr, err := buildResponse(ext, ext.SequenceResponse.AfterEnded)
			if err != nil {
				return nil, err
			}

			seq.AfterEnded(rr)
		}

		for _, r := range ext.SequenceResponse.Responses {
			rr, err := buildResponse(ext, &r)
			if err != nil {
				return nil, err
			}

			seq.Add(rr)
		}

		rep = seq
	} else {
		rep = reply.OK()
	}

	builder.Reply(rep)

	return builder, nil
}

func buildResponse(ext *mod.ExternalSchema, response *mod.ExtRes) (reply.Reply, error) {
	res := reply.New()
	res.Status(valueOr(response.Status, http.StatusOK))

	for k, v := range response.Header {
		res.Header(k, v)
	}

	if response.BodyFile != "" {
		file, err := os.Open(response.BodyFile)
		if err != nil {
			return nil, err
		}

		defer file.Close()

		b, err := io.ReadAll(file)
		if err != nil {
			return nil, err
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
