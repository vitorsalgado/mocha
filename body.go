package mocha

import (
	"encoding/json"
	"net/http"
	"strings"
)

type (
	BodyParser interface {
		CanParse(content string, r *http.Request) bool
		Parse(r *http.Request, v any) error
	}

	Parsers struct {
		parsers []BodyParser
	}
)

func (p *Parsers) Compose(parsers ...BodyParser) {
	p.parsers = append(p.parsers, parsers...)
}

func (p *Parsers) Get() []BodyParser {
	return p.parsers
}

type JSONBodyParser struct{}

func (parser JSONBodyParser) CanParse(content string, _ *http.Request) bool {
	return strings.Contains(content, ContentTypeJSON)
}

func (parser JSONBodyParser) Parse(r *http.Request, v any) error {
	return json.NewDecoder(r.Body).Decode(v)
}

type FormURLEncodedParser struct{}

func (parser FormURLEncodedParser) CanParse(content string, _ *http.Request) bool {
	return strings.Contains(content, ContentTypeFormURLEncoded)
}

func (parser *FormURLEncodedParser) Parse(r *http.Request, v any) error {
	err := r.ParseForm()
	if err != nil {
		return err
	}

	v = r.Form.Encode()

	return nil
}
