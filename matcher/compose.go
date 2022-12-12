package matcher

type ComposeMatcher struct {
	M Matcher
}

func (m *ComposeMatcher) Name() string                 { return m.M.Name() }
func (m *ComposeMatcher) Match(v any) (*Result, error) { return m.M.Match(v) }
func (m *ComposeMatcher) OnMockServed() error          { return m.M.OnMockServed() }
func (m *ComposeMatcher) Spec() any                    { return m.Spec() }

// And compose the current Matcher with another one using the "and" operator.
func (m *ComposeMatcher) And(and Matcher) *ComposeMatcher {
	return Compose(AllOf(m, and))
}

// Or compose the current Matcher with another one using the "or" operator.
func (m *ComposeMatcher) Or(or Matcher) *ComposeMatcher {
	return Compose(AnyOf(m, or))
}

// Xor compose the current Matcher with another one using the "xor" operator.
func (m *ComposeMatcher) Xor(and Matcher) *ComposeMatcher {
	return Compose(XOR(m, and))
}

func Compose(base Matcher) *ComposeMatcher {
	return &ComposeMatcher{M: base}
}
