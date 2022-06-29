package params

type (
	Params struct {
		param map[string]any
	}
)

func New() *Params {
	return &Params{param: make(map[string]any)}
}

func (p Params) Get(key string) (any, bool) {
	val, ok := p.param[key]
	return val, ok
}

func (p Params) GetAll() map[string]any {
	return p.param
}

func (p Params) Set(key string, dep any) {
	p.param[key] = dep
}

func (p Params) Remove(key string) {
	delete(p.param, key)
}

func (p Params) Has(key string) bool {
	_, ok := p.param[key]

	return ok
}
