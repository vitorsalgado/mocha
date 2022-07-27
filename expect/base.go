package expect

import "net/http"

type (
	// Params define a custom parameter holder for matchers.
	Params interface {
		Get(key string) (any, bool)
	}

	// RequestInfo implements HTTP request information to be passed to each Matcher.
	RequestInfo struct {
		// Request is the actual http.Request.
		Request *http.Request

		// ParsedBody is http.Request parsed body.
		// Value of parsed body can vary depending on the mocha.RequestBodyParser that parsed the request.
		ParsedBody any
	}

	// Args groups contextual information available for each Matcher.
	Args struct {
		RequestInfo *RequestInfo
		Params      Params
	}

	// Matcher defines request matchers.
	// Request matchers are used to match requests in order to find a mock to serve a stub response.
	Matcher struct {
		// Name is a metadata that defines the matcher name.
		Name string

		// DescribeMismatch gives more context of why the Matcher failed to match a given value.
		DescribeMismatch func(p string, v any) string

		// Matches is the function that does the actual matching logic.
		Matches func(v any, args Args) (bool, error)
	}
)

// And compose the current Matcher with another one using the "and" operator.
func (m Matcher) And(and Matcher) Matcher {
	return AllOf(m, and)
}

// Or compose the current Matcher with another one using the "or" operator.
func (m Matcher) Or(or Matcher) Matcher {
	return AnyOf(m, or)
}

// Xor compose the current Matcher with another one using the "xor" operator.
func (m Matcher) Xor(and Matcher) Matcher {
	return XOR(m, and)
}
