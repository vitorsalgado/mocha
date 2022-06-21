package mocha

import (
	"net/http"

	"github.com/vitorsalgado/mocha/matcher"
)

func WrapRequest(r *http.Request, parsers []BodyParser) (*matcher.RequestInfo, error) {
	if r.Body != nil && r.Method != http.MethodGet {
		var content = r.Header.Get("content-type")

		for _, parse := range parsers {
			if parse.CanParse(content, r) {
				body, err := parse.Parse(r)
				if err != nil {
					return nil, err
				}

				return &matcher.RequestInfo{Request: r, Body: body}, nil
			}
		}
	}

	return &matcher.RequestInfo{Request: r}, nil
}
