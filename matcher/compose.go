package matcher

type ComposeMatcher struct {
	M Matcher
}

func (m *ComposeMatcher) Name() string                 { return m.M.Name() }
func (m *ComposeMatcher) Match(v any) (*Result, error) { return m.M.Match(v) }
func (m *ComposeMatcher) AfterMockServed() error       { return runAfterMockServed(m.M) }

// And compose the current Matcher with another one using the "and" operator.
func (m *ComposeMatcher) And(and Matcher) *ComposeMatcher {
	return Compose(All(m, and))
}

// Or compose the current Matcher with another one using the "or" operator.
func (m *ComposeMatcher) Or(or Matcher) *ComposeMatcher {
	return Compose(Any(m, or))
}

// Xor compose the current Matcher with another one using the "xor" operator.
func (m *ComposeMatcher) Xor(and Matcher) *ComposeMatcher {
	return Compose(XOR(m, and))
}

func Compose(base Matcher) *ComposeMatcher {
	return &ComposeMatcher{M: base}
}
