package mocha

import (
	"net/http"
)

type (
	MockRequest struct {
		RawRequest *http.Request
		Body       any
	}
)

func WrapRequest(r *http.Request, parsers []BodyParser) (*MockRequest, error) {
	if r.Body != nil && r.Method != http.MethodGet {
		var content = r.Header.Get("content-type")

		for _, parse := range parsers {
			if parse.CanParse(content, r) {
				body, err := parse.Parse(r)
				if err != nil {
					return nil, err
				}

				return &MockRequest{RawRequest: r, Body: body}, nil
			}
		}
	}

	return &MockRequest{RawRequest: r}, nil
}
