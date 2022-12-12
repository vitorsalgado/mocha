package mocha

type target int

// target constants to help debug unmatched requests.
const (
	_targetRequest target = iota
	_targetMethod
	_targetURL
	_targetHeader
	_targetQuery
	_targetBody
	_targetForm
)
