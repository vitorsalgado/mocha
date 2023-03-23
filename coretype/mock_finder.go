package coretype

import (
	"reflect"
)

// FindResult holds the results for an attempt to match a mock to a request.
type FindResult[MOCK Mock] struct {
	Pass            bool
	Matched         MOCK
	ClosestMatch    MOCK
	MismatchDetails []MismatchDetail
}

type MMO[TValueIn any] interface {
	Mock
	RequestMatcher[TValueIn]
}

// FindMockForRequest tries to find a mock to the incoming HTTP request.
// It runs all matchers of all eligible mocks on request until it finds one that matches every one of them.
// It returns a FindResult with the find result, along with a possible closest match.
func FindMockForRequest[TValueIn any, MOCK MMO[TValueIn]](
	storage *MockStore[MOCK],
	requestValues TValueIn,
) *FindResult[MOCK] {
	var mocks = storage.GetEligible()
	var matched MOCK
	var weights = 0
	var details = make([]MismatchDetail, 0)

	for _, m := range mocks {
		result := Match[TValueIn](requestValues, m.GetExpectations())

		if result.Pass {
			return &FindResult[MOCK]{Pass: true, Matched: m}
		}

		if result.Weight > 0 && result.Weight > weights {
			matched = m
			weights = result.Weight
		}

		details = append(details, result.Details...)
	}

	if reflect.ValueOf(matched).IsNil() {
		return &FindResult[MOCK]{Pass: false, MismatchDetails: details}
	}

	return &FindResult[MOCK]{Pass: false, ClosestMatch: matched, MismatchDetails: details}
}
