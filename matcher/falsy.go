package matcher

type falsyMatcher struct {
}

func (m *falsyMatcher) Name() string {
	return "Falsy"
}

func (m *falsyMatcher) Match(v any) (*Result, error) {
	res, err := Truthy().Match(v)
	if err != nil {
		return nil, err
	}

	return &Result{Pass: !res.Pass, Ext: res.Ext, Message: res.Message}, nil
}

func Falsy() Matcher {
	return &falsyMatcher{}
}
