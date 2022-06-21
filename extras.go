package mocha

type Extras struct {
	deps map[string]any
}

func NewExtras() Extras {
	return Extras{deps: make(map[string]any)}
}

func (deps *Extras) Get(key string) (any, bool) {
	val, ok := deps.deps[key]
	return val, ok
}

func (deps *Extras) Set(key string, dep any) {
	deps.deps[key] = dep
}

func (deps *Extras) Remove(key string) {
	delete(deps.deps, key)
}
