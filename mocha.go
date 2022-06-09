package mocha

import "net/http/httptest"

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
	return &Mocha{Server: httptest.NewServer(&Handler{repo: repo}), Repo: repo}
}

func (m *Mocha) Start() Info {
	m.Server.Start()
	return Info{URL: m.Server.URL}
}

func (m *Mocha) Mock(mock Mock) *Mocha {
	m.Repo.Save(mock)
	return m
}

func (m *Mocha) Close() {
	m.Server.Close()
}
