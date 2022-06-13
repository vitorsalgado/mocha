package mocha

import (
	"net/http/httptest"
	"testing"
)

type (
	Mocha struct {
		Server *httptest.Server
		Repo   MockRepository
	}

	Info struct {
		URL string
	}
)

func New() *Mocha {
	repo := NewMockRepository()
	sp := NewScenarioRepository()
	return &Mocha{Server: httptest.NewServer(&Handler{repo: repo, scenarioRepository: sp}), Repo: repo}
}

func NewT(t *testing.T) *Mocha {
	m := New()
	t.Cleanup(func() { m.Close() })

	return m
}

func (m *Mocha) Start() Info {
	m.Server.Start()
	return Info{URL: m.Server.URL}
}

func (m *Mocha) Mock(builders ...*MockBuilder) *Scoped {
	added := make([]int32, 0)

	for _, b := range builders {
		mock := b.Build()

		m.Repo.Save(mock)
		added = append(added, mock.ID)
	}

	return NewScoped(m.Repo, added)
}

func (m *Mocha) Close() {
	m.Server.Close()
}
