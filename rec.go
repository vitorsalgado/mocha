package mocha

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"mime"
	"net/http"
	"net/url"
	"os"
	"path"
	"regexp"
	"strings"
	"sync"

	"github.com/vitorsalgado/mocha/v3/internal/header"
	"github.com/vitorsalgado/mocha/v3/mod"
)

var _regexNoSpecialCharacters = regexp.MustCompile("[^a-z0-9]")
var _defaultRecordConfig = RecordConfig{
	SaveDir:           "testdata",
	Save:              false,
	RecRequestHeaders: []string{"accept", "content-type"},
	RecResponseHeaders: []string{
		"content-type",
		"link",
		"content-length",
		"cache-control",
		"retry-after"},
}

type RecordConfig struct {
	RecRequestHeaders  []string
	RecResponseHeaders []string
	Save               bool
	SaveDir            string
}

type RecordConfigurer interface {
	Apply(c *RecordConfig)
}

func (r *RecordConfig) Apply(opts *RecordConfig) {
	opts.SaveDir = r.SaveDir
	opts.RecResponseHeaders = r.RecResponseHeaders
	opts.RecRequestHeaders = r.RecRequestHeaders
}

type recArgs struct {
	request  recRequest
	response recResponse
}

type recRequest struct {
	path   string
	method string
	header http.Header
	query  url.Values
	body   []byte
}

type recResponse struct {
	status int
	header http.Header
	body   []byte
}

type record struct {
	config *RecordConfig
	in     chan *recArgs
	cancel context.CancelFunc
	mu     sync.Mutex
	data   []*Mock
}

func newRecord(config *RecordConfig) *record {
	return &record{config: config}
}

func (r *record) startRecording(ctx context.Context) {
	ctx, cancel := context.WithCancel(ctx)

	r.in = make(chan *recArgs)
	r.cancel = cancel

	go func() {
		for {
			select {
			case a, ok := <-r.in:
				if !ok {
					return
				}

				r.process(a)

			case <-ctx.Done():
				close(r.in)
				return
			}
		}
	}()
}

func (r *record) record(arg *recArgs) {
	r.in <- arg
}

func (r *record) stop() {
	r.cancel()
}

func (r *record) process(arg *recArgs) {
	out := &mod.ExtMock{}
	out.Response = &mod.ExtMockResponse{}
	out.Response.Header = make(map[string]string)

	sanitized := strings.ReplaceAll(arg.request.path, " ", "-")
	sanitized = strings.TrimPrefix(_regexNoSpecialCharacters.ReplaceAllString(sanitized, "-"), "-")
	name := fmt.Sprintf("%s-%s", arg.request.method, sanitized)

	out.Request.URL = arg.request.path
	out.Request.Method = arg.request.method

	if len(r.config.RecRequestHeaders) > 0 {
		for _, h := range r.config.RecRequestHeaders {
			v := arg.request.header.Get(h)
			if v != "" {
				out.Request.Header[h] = v
			}
		}
	}

	if len(r.config.RecResponseHeaders) > 0 {
		for _, h := range r.config.RecResponseHeaders {
			v := arg.response.header.Get(h)
			if v != "" {
				out.Response.Header[h] = v
			}
		}
	}

	for k := range arg.request.query {
		out.Request.Query[k] = arg.request.query.Get(k)
	}

	if len(arg.request.body) > 0 {
		out.Request.Body = []any{"equal", string(arg.request.body)}
	}

	out.Response.Status = arg.response.status

	contentType := arg.request.header.Get(header.ContentType)
	if strings.Contains(contentType, ";") {
		contentType = strings.TrimSpace(contentType[:strings.Index(contentType, ";")])
	}
	ext, err := mime.ExtensionsByType(contentType)
	if err != nil {
		log.Print(err)
	}

	dst := r.config.SaveDir
	if !path.IsAbs(r.config.SaveDir) {
		wd, _ := os.Getwd()
		dst = path.Join(wd, dst)
	}

	mockFile := path.Join(dst, fmt.Sprintf("%s.mock.json", name))
	mockBodyFile := ""

	_, err = os.Stat(mockFile)
	exists := true
	if err != nil {
		exists = false
	}

	if exists {
		log.Printf("file %s exists.\n", mockFile)
		return
	}

	if len(arg.response.body) > 0 {
		if len(ext) > 0 {
			mockBodyFile = name + ext[0]
		} else {
			mockBodyFile = name + ".bin"
		}

		b, err := os.Create(path.Join(dst, mockBodyFile))
		if err != nil {
			log.Println(err)
			return
		}

		defer b.Close()

		_, err = b.Write(arg.response.body)
		if err != nil {
			log.Println(err)
			return
		}

		out.Response.BodyFile = mockBodyFile
	}

	file, err := os.Create(mockFile)
	if err != nil {
		log.Println(err)
		return
	}

	defer file.Close()

	err = json.NewEncoder(file).Encode(out)
	if err != nil {
		log.Println(err)
		return
	}
}

type recordConfigFunc func(config *RecordConfig)

func (f recordConfigFunc) Apply(config *RecordConfig) { f(config) }

func WithRecordDir(destination string) RecordConfigurer {
	return recordConfigFunc(func(c *RecordConfig) { c.SaveDir = destination })
}
