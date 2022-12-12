package mocha

import (
	"context"
	"fmt"
	"path/filepath"
	"sync"
)

// Loader is the interface that defines custom Mock loaders.
// Usually, it is used to load external mocks, like from the file system.
type Loader interface {
	Load(app *Mocha) error
}

var _ Loader = (*FileLoader)(nil)

type FileLoader struct {
	mu sync.Mutex
}

func (l *FileLoader) Load(app *Mocha) error {
	l.mu.Lock()
	defer l.mu.Unlock()

	filenames := make(map[string]struct{})

	for _, pattern := range app.Config.Directories {
		m, err := filepath.Glob(pattern)
		if err != nil {
			return fmt.Errorf("error searching mocks with the glob pattern %s. %w", pattern, err)
		}

		for _, s := range m {
			if _, ok := filenames[s]; !ok {
				filenames[s] = struct{}{}
			}
		}
	}

	type errContainer struct {
		err error
	}

	cont := &errContainer{}
	ch := make(chan string, len(filenames))
	wg := sync.WaitGroup{}
	once := sync.Once{}
	w := 5
	ctx, cancel := context.WithCancel(app.Context())

	if len(filenames) < w {
		w = len(filenames)
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
						_, err := app.Mock(From(filename))
						return err
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

	for filename := range filenames {
		wg.Add(1)
		ch <- filename
	}

	wg.Wait()
	cancel()
	close(ch)

	return cont.err
}
