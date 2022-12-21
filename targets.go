package mocha

type matchTarget int

// matchTarget constants to help debug unmatched requests.
const (
	_targetRequest matchTarget = iota
	_targetMethod
	_targetURL
	_targetHeader
	_targetQuery
	_targetBody
	_targetForm
)
