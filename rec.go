package mocha

import (
	"context"
	"encoding/hex"
	"fmt"
	"hash/fnv"
	"log"
	"mime"
	"net/http"
	"net/url"
	"os"
	"path"
	"regexp"
	"strings"

	"github.com/spf13/viper"

	"github.com/vitorsalgado/mocha/v3/internal/header"
	"github.com/vitorsalgado/mocha/v3/reply"
)

var (
	_regexNoSpecialCharacters = regexp.MustCompile("[^a-z0-9]")
)

type RecordConfig struct {
	RequestHeaders  []string
	ResponseHeaders []string
	Save            bool
	SaveDir         string
	SaveExtension   string
	SaveBodyToFile  bool
}

func defaultRecordConfig() *RecordConfig {
	return &RecordConfig{
		SaveDir:        "testdata/_mocks",
		SaveExtension:  "json",
		SaveBodyToFile: false,
		Save:           false,
		RequestHeaders: []string{"accept", "content-type"},
		ResponseHeaders: []string{
			"content-type",
			"link",
			"content-length",
			"cache-control",
			"retry-after"},
	}
}

type RecordConfigurer interface {
	Apply(c *RecordConfig)
}

func (r *RecordConfig) Apply(opts *RecordConfig) {
	opts.Save = r.Save
	opts.SaveDir = r.SaveDir
	opts.SaveBodyToFile = r.SaveBodyToFile
	opts.ResponseHeaders = r.ResponseHeaders
	opts.RequestHeaders = r.RequestHeaders
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

				err := r.process(a)
				if err != nil {
					log.Println(err)
				}

			case <-ctx.Done():
				close(r.in)
				return
			}
		}
	}()
}

func (r *record) record(req *http.Request, rawReqBody []byte, res *reply.Stub) {
	input := &recArgs{
		request: recRequest{
			path:   req.URL.Path,
			method: req.Method,
			header: req.Header,
			query:  req.URL.Query(),
			body:   rawReqBody,
		},
		response: recResponse{
			status: res.StatusCode,
			header: res.Header,
			body:   res.Body,
		},
	}

	go func() { r.in <- input }()
}

func (r *record) stop() {
	r.cancel()
}

func (r *record) process(arg *recArgs) error {
	v := viper.New()
	v.AddConfigPath(r.config.SaveDir)
	v.SetConfigType(r.config.SaveExtension)

	v.Set(_fRequestURLPath, arg.request.path)
	v.Set(_fRequestMethod, arg.request.method)
	v.Set(_fResponseStatus, arg.response.status)

	requestHeaders := make(map[string]string)
	responseHeaders := make(map[string]string)
	query := make(map[string]string, len(arg.request.query))
	hasResBody := len(arg.response.body) > 0
	bsha := ""

	for _, h := range r.config.RequestHeaders {
		v := arg.request.header.Get(h)
		if v != "" {
			requestHeaders[h] = v
		}
	}

	for _, h := range r.config.ResponseHeaders {
		v := arg.response.header.Get(h)
		if v != "" {
			responseHeaders[h] = v
		}
	}

	for k := range arg.request.query {
		query[k] = arg.request.query.Get(k)
	}

	if len(arg.request.body) > 0 {
		v.Set(_fRequestBody, []string{"equal", string(arg.request.body)})

		h := fnv.New64()
		h.Write(arg.request.body)
		bsha = hex.EncodeToString(h.Sum(nil))
	}

	contentType := arg.request.header.Get(header.ContentType)
	if strings.Contains(contentType, ";") {
		contentType = strings.TrimSpace(contentType[:strings.Index(contentType, ";")])
	}
	ext, _ := mime.ExtensionsByType(contentType)

	nm := strings.ReplaceAll(arg.request.path, " ", "-")
	nm = strings.TrimPrefix(_regexNoSpecialCharacters.ReplaceAllString(nm, "-"), "-")
	if bsha != "" {
		nm += "--" + bsha
	}

	name := fmt.Sprintf("%s-%s", arg.request.method, nm)
	mockFile := path.Join(r.config.SaveDir, fmt.Sprintf("%s.mock.json", name))

	_, err := os.Stat(mockFile)
	exists := true
	if err != nil {
		exists = false
	}

	if exists {
		return fmt.Errorf("file %s exists", mockFile)
	}

	if hasResBody {
		if r.config.SaveBodyToFile {
			mockBodyFile := name + ".bin"
			if len(ext) > 0 {
				mockBodyFile = name + ext[0]
			}

			b, err := os.Create(path.Join(r.config.SaveDir, mockBodyFile))
			if err != nil {
				return err
			}

			defer b.Close()

			_, err = b.Write(arg.response.body)
			if err != nil {
				return err
			}

			v.Set(_fResponseBodyFile, mockBodyFile)
		} else {
			v.Set(_fResponseBody, string(arg.response.body))
		}
	}

	if r.config.Save {
		err = v.WriteConfigAs(mockFile)
		if err != nil {
			return err
		}
	}

	return nil
}

type recordConfigFunc func(config *RecordConfig)

func (f recordConfigFunc) Apply(config *RecordConfig) { f(config) }

func WithRecordDir(destination string) RecordConfigurer {
	return recordConfigFunc(func(c *RecordConfig) { c.SaveDir = destination })
}
