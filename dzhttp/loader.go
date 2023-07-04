package dzhttp

import (
	"context"
	"fmt"
	"path/filepath"
	"sync"
)

// Loader is the interface that defines custom Mock loaders.
// Usually, it is used to load external mocks, like from the file system.
type Loader interface {
	Load(app *HTTPMockApp) error
}

var _ Loader = (*fileLoader)(nil)

type fileLoader struct {
	rwMutex sync.RWMutex
}

func (l *fileLoader) Load(app *HTTPMockApp) error {
	l.rwMutex.Lock()
	defer l.rwMutex.Unlock()

	filenames := make(map[string]struct{})

	for _, pattern := range app.config.MockFileSearchPatterns {
		m, err := filepath.Glob(filepath.Join(app.config.RootDir, pattern))
		if err != nil {
			return fmt.Errorf("loader: error searching mocks with the glob pattern %s.\n%w", pattern, err)
		}

		for _, s := range m {
			if _, ok := filenames[s]; !ok {
				filenames[s] = struct{}{}
			}
		}
	}

	if len(filenames) == 0 {
		return nil
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
						_, err := app.Mock(FromFile(filename))
						if err != nil {
							return fmt.Errorf("loader: error adding mock.\n%w", err)
						}

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

	for filename := range filenames {
		wg.Add(1)
		ch <- filename
	}

	wg.Wait()
	cancel()
	close(ch)

	return cont.err
}
