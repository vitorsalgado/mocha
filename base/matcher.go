package base

type MatcherContext struct {
}

type Matcher[V any] func(v V, ctx MatcherContext) (bool, error)
