package matcher

import (
	"fmt"
	"reflect"
)

type itemsMatchMatcher struct {
	expected any
}

func (m *itemsMatchMatcher) Name() string {
	return "EqualElements"
}

func (m *itemsMatchMatcher) Match(v any) (*Result, error) {
	a := reflect.ValueOf(v)
	b := reflect.ValueOf(m.expected)

	aLen := a.Len()
	bLen := b.Len()

	if aLen != bLen {
		return &Result{
			Ext: []string{stringify(m.expected)},
			Message: fmt.Sprintf(
				"Expected value length %d. Received length %d",
				aLen,
				bLen)}, nil
	}

	var extraA, extraB []interface{}

	visited := make([]bool, bLen)

	for i := 0; i < aLen; i++ {
		element := a.Index(i).Interface()
		found := false

		for j := 0; j < bLen; j++ {
			if visited[j] {
				continue
			}

			if equalValues(b.Index(j).Interface(), element) {
				visited[j] = true
				found = true
				break
			}
		}

		if !found {
			extraA = append(extraA, element)
		}
	}

	for j := 0; j < bLen; j++ {
		if visited[j] {
			continue
		}

		extraB = append(extraB, b.Index(j).Interface())
	}

	if len(extraA) == 0 && len(extraB) == 0 {
		return &Result{Pass: true}, nil
	}

	return &Result{}, nil
}

func ItemsMatch(items any) Matcher {
	return &itemsMatchMatcher{expected: items}
}
