package params

type (
	Params interface {
		Get(key string) (any, bool)
		GetAll() map[string]any
		Set(key string, dep any)
		Remove(key string)
		Has(key string) bool
	}

	params struct {
		data map[string]any
	}
)

func New() Params {
	return &params{data: make(map[string]any)}
}

func (p params) Get(key string) (any, bool) {
	val, ok := p.data[key]
	return val, ok
}

func (p params) GetAll() map[string]any {
	return p.data
}

func (p params) Set(key string, dep any) {
	p.data[key] = dep
}

func (p params) Remove(key string) {
	delete(p.data, key)
}

func (p params) Has(key string) bool {
	_, ok := p.data[key]

	return ok
}
