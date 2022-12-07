package mocha

import (
	"context"
	"encoding/json"
	"fmt"
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
			return fmt.Errorf("error searching mocks with the glob pattern %s. %w", pattern, err)
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

	cont := &errContainer{}
	ch := make(chan string, len(matches))
	wg := sync.WaitGroup{}
	once := sync.Once{}
	w := 5
	ctx, cancel := context.WithCancel(app.Context())

	if len(matches) < 5 {
		w = len(matches)
	}

	for i := 0; i < w; i++ {
		go func(c *errContainer) {
			for {
				select {
				case filename, ok := <-ch:
					if !ok {
						return
					}

					fn := func(filename string, c *errContainer) error {
						file, err := os.Open(filename)
						if err != nil {
							return fmt.Errorf("error opening mock file [%s]. %w", filename, err)
						}

						v, err := l.decode(filename, file)
						if err != nil {
							return fmt.Errorf("error decoding mock [%s]. %w", filename, err)
						}

						defer file.Close()

						m, err := buildExternalMock(filename, v)
						if err != nil {
							return fmt.Errorf("error building mock [%s]. %w", filename, err)
						}

						app.AddMocks(m)

						return nil
					}

					err := fn(filename, c)
					if err != nil {
						once.Do(func() {
							c.err = err
							wg.Done()
							cancel()
						})
					} else {
						wg.Done()
					}

				case <-ctx.Done():
					return
				}
			}
		}(cont)
	}

	for _, filename := range matches {
		wg.Add(1)
		ch <- filename
	}

	wg.Wait()
	cancel()
	close(ch)

	return cont.err
}

func (l *FileLoader) decode(filename string, file io.ReadCloser) (r *mod.ExtMock, err error) {
	parts := strings.Split(filename, "/")
	ext := parts[len(parts)-1]

	switch ext {
	// JSON is default
	default:
		r := &mod.ExtMock{}
		err = json.NewDecoder(file).Decode(&r)
		return r, err
	}
}
