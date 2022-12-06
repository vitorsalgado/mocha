package mocha

import (
	"context"
	"encoding/json"
	"io"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/vitorsalgado/mocha/v3/mod"
)

var _ Loader = (*FileLoader)(nil)

type FileLoader struct {
}

type errContainer struct {
	err error
}

func (l *FileLoader) Load(app *Mocha) error {
	set := make(map[string]struct{})

	for _, pattern := range app.Config.FileMockPatterns {
		m, err := filepath.Glob(pattern)
		if err != nil {
			return err
		}

		for _, s := range m {
			if _, ok := set[s]; !ok {
				set[s] = struct{}{}
			}
		}
	}

	matches := make([]string, len(set))
	i := 0
	for k := range set {
		matches[i] = k
		i++
	}

	c := &errContainer{}
	wg := sync.WaitGroup{}
	once := sync.Once{}
	_, cancel := context.WithCancel(context.Background())
	defer cancel()

	for _, filename := range matches {
		wg.Add(1)

		go func(filename string, c *errContainer) {
			defer wg.Done()

			fn := func(filename string, c *errContainer) error {
				file, err := os.Open(filename)
				if err != nil {
					app.T.Logf("error loading mock file %s. reason=%s", filename, err.Error())
					return err
				}

				v, err := l.decode(filename, file)
				if err != nil {
					app.T.Logf("error decoding mock file %s. reason=%s", filename, err.Error())
					return err
				}

				defer file.Close()

				m, err := buildExternalMock(filename, v)
				if err != nil {
					app.T.Logf("error building mock file %s. reason=%s", filename, err.Error())
					return err
				}

				app.AddMocks(m)

				return nil
			}

			err := fn(filename, c)
			if err != nil {
				once.Do(func() {
					c.err = err
					cancel()
				})
			}
		}(filename, c)
	}

	wg.Wait()

	return c.err
}

func (l *FileLoader) decode(filename string, file io.ReadCloser) (r *mod.ExternalSchema, err error) {
	parts := strings.Split(filename, "/")
	ext := parts[len(parts)-1]

	switch ext {
	// JSON is default
	default:
		r := &mod.ExternalSchema{}
		err = json.NewDecoder(file).Decode(&r)
		return r, err
	}
}
