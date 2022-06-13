package mocha

import (
	"net/http"
)

type MatcherContext struct {
	Config             Config
	Req                *http.Request
	Repo               MockRepository
	ScenarioRepository ScenarioRepository
}

type Matcher[V any] func(v V, ctx MatcherContext) (bool, error)
