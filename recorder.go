package mocha

import (
	"bytes"
	"compress/gzip"
	"context"
	"encoding/hex"
	"fmt"
	"hash/fnv"
	"io"
	"log"
	"mime"
	"net/http"
	"net/url"
	"os"
	"path"
	"regexp"
	"strings"
	"sync"

	"github.com/spf13/viper"

	"github.com/vitorsalgado/mocha/v3/internal/header"
)

var _ RecordConfigurer = (*RecordConfig)(nil)

var _regexNoSpecialCharacters = regexp.MustCompile("[^a-z0-9]")

// RecordConfig configures HTTP request and response recording.
type RecordConfig struct {
	// RequestHeaders define the request headers that should be recorded.
	// Defaults to: "accept", "content-type"
	RequestHeaders []string

	// ResponseHeaders define the response headers that should be recorded.
	// Defaults to:
	//	"content-type"
	//	"link"
	//	"content-length"
	//	"cache-control"
	//	"retry-after"
	ResponseHeaders []string

	// SaveDir defines the directory to save the recorded mocks.
	// Defaults to: testdata/mocks
	SaveDir string

	// SaveExtension defines the extension to save the recorded mock.
	// It uses viper.Viper to save the mock files, so it will accept the same values that viper supports.
	// Defaults to .json.
	SaveExtension string

	// SaveResponseBodyToFile defines if the recorded body should be saved to a separate file.
	// Defaults to false, embedding the response body into the mock definition file.
	SaveResponseBodyToFile bool
}

// RecordConfigurer lets users configure recording.
type RecordConfigurer interface {
	Apply(c *RecordConfig) error
}

type recordConfigFunc func(config *RecordConfig)

func (f recordConfigFunc) Apply(config *RecordConfig) error {
	f(config)
	return nil
}

// Apply allows using RecordConfig as RecordConfigurer.
func (r *RecordConfig) Apply(opts *RecordConfig) error {
	extension := strings.TrimPrefix(opts.SaveExtension, ".")
	found := false
	for _, ext := range viper.SupportedExts {
		if strings.EqualFold(extension, ext) {
			found = true
			break
		}
	}

	if !found {
		return fmt.Errorf(
			"recorder: recorded mock file extension \"%s\" is not supported. supported values: %s",
			extension,
			strings.Join(viper.SupportedExts, ", "),
		)
	}

	opts.SaveExtension = extension
	opts.SaveDir = r.SaveDir
	opts.SaveResponseBodyToFile = r.SaveResponseBodyToFile
	opts.ResponseHeaders = r.ResponseHeaders
	opts.RequestHeaders = r.RequestHeaders

	return nil
}

func defaultRecordConfig() *RecordConfig {
	return &RecordConfig{
		SaveDir:                "testdata/mocks",
		SaveExtension:          "json",
		SaveResponseBodyToFile: false,
		RequestHeaders:         []string{"accept", "content-type", "content-encoding", "content-length"},
		ResponseHeaders:        []string{"content-type", "content-encoding", "link", "content-length", "cache-control", "retry-after"},
	}
}

// Recording data transfer structs.
type (
	recArgs struct {
		request  recRequest
		response recResponse
	}

	recRequest struct {
		uri    string
		method string
		header http.Header
		query  url.Values
		body   []byte
	}

	recResponse struct {
		status int
		header http.Header
		body   []byte
	}

	recorder struct {
		active bool
		config *RecordConfig
		in     chan *recArgs
		cancel context.CancelFunc
		mu     sync.Mutex
	}
)

func newRecorder(config *RecordConfig) *recorder {
	return &recorder{config: config}
}

func (r *recorder) start(ctx context.Context) {
	r.mu.Lock()
	defer r.mu.Unlock()

	ctx, cancel := context.WithCancel(ctx)

	r.in = make(chan *recArgs)
	r.cancel = cancel
	r.active = true

	go func() {
		for {
			select {
			case a, ok := <-r.in:
				if !ok {
					return
				}

				err := r.process(a)
				if err != nil {
					log.Println(fmt.Errorf("recorder: %w", err))
				}

			case <-ctx.Done():
				return
			}
		}
	}()
}

func (r *recorder) dispatch(req *http.Request, parsedURL *url.URL, rawReqBody []byte, res *Stub) {
	if !r.active {
		return
	}

	input := &recArgs{
		request: recRequest{
			uri:    parsedURL.RequestURI(),
			method: req.Method,
			header: req.Header.Clone(),
			query:  req.URL.Query(),
			body:   rawReqBody,
		},
		response: recResponse{
			status: res.StatusCode,
			header: res.Header.Clone(),
			body:   res.Body,
		},
	}

	go func() { r.in <- input }()
}

func (r *recorder) stop() {
	r.active = false
	r.cancel()
}

func (r *recorder) process(arg *recArgs) error {
	v := viper.New()

	v.AddConfigPath(r.config.SaveDir)
	v.SetConfigType(r.config.SaveExtension)

	v.Set(_fRequestURLPath, arg.request.uri)
	v.Set(_fRequestMethod, arg.request.method)

	requestHeaders := make(map[string]string)
	responseHeaders := make(map[string]string)
	query := make(map[string]string, len(arg.request.query))
	hasResBody := len(arg.response.body) > 0
	requestBodyHash := ""

	for _, h := range r.config.RequestHeaders {
		headerValue := arg.request.header.Get(h)
		if headerValue != "" {
			requestHeaders[h] = headerValue
		}
	}

	for _, h := range r.config.ResponseHeaders {
		headerValue := arg.response.header.Get(h)
		if headerValue != "" {
			responseHeaders[h] = headerValue
		}
	}

	for k := range arg.request.query {
		query[k] = arg.request.query.Get(k)
	}

	v.Set(_fRequestQuery, query)
	v.Set(_fRequestHeader, requestHeaders)

	if len(arg.request.body) > 0 {
		v.Set(_fRequestBody, []string{"equal", string(arg.request.body)})

		h := fnv.New64()
		h.Write(arg.request.body)
		requestBodyHash = hex.EncodeToString(h.Sum(nil))
	}

	contentType := arg.request.header.Get(header.ContentType)
	if strings.Contains(contentType, ";") {
		contentType = strings.TrimSpace(contentType[:strings.Index(contentType, ";")])
	}

	nm := strings.ReplaceAll(arg.request.uri, " ", "-")
	nm = strings.TrimPrefix(_regexNoSpecialCharacters.ReplaceAllString(nm, "-"), "-")
	if requestBodyHash != "" {
		nm += "--" + requestBodyHash
	}

	name := fmt.Sprintf("%s-%s", arg.request.method, nm)
	mockFile := path.Join(r.config.SaveDir, fmt.Sprintf("%s.mock.%s", name, r.config.SaveExtension))

	_, err := os.Stat(mockFile)
	exists := true
	if err != nil {
		exists = false
	}

	if exists {
		return fmt.Errorf("file %s already exists", mockFile)
	}

	if hasResBody {
		if r.config.SaveResponseBodyToFile {
			bodyFilename := name + "--response-body"
			ext, _ := mime.ExtensionsByType(contentType)

			if len(ext) > 0 {
				bodyFilename += ext[0]
			} else {
				bodyFilename += ".bin"
			}

			b, err := os.Create(path.Join(r.config.SaveDir, bodyFilename))
			if err != nil {
				return err
			}

			defer b.Close()

			encoding := arg.response.header.Get(header.ContentEncoding)
			switch encoding {
			case "gzip":
				gz, err := gzip.NewReader(bytes.NewReader(arg.response.body))
				if err != nil {
					return err
				}

				defer gz.Close()

				body, err := io.ReadAll(gz)
				if err != nil {
					return err
				}

				b.Write(body)
				v.Set(_fResponseEncoding, "gzip")
			default:
				b.Write(arg.response.body)
			}

			v.Set(_fResponseBodyFile, bodyFilename)

		} else {
			encoding := arg.response.header.Get(header.ContentEncoding)
			switch encoding {
			case "gzip":
				gz, err := gzip.NewReader(bytes.NewReader(arg.response.body))
				if err != nil {
					return err
				}

				defer gz.Close()

				b, err := io.ReadAll(gz)
				if err != nil {
					return err
				}

				v.Set(_fResponseBody, string(b))
				v.Set(_fResponseEncoding, "gzip")
			default:
				v.Set(_fResponseBody, string(arg.response.body))
			}
		}
	}

	if v.GetBool(_fResponseTemplateEnabled) {
		v.Set(_fResponseTemplateEnabled, v.GetBool(_fResponseTemplateEnabled))
		v.Set(_fResponseTemplateModel, v.Get(_fResponseTemplateModel))
	}

	v.Set(_fResponseStatus, arg.response.status)
	v.Set(_fResponseHeader, responseHeaders)

	err = v.WriteConfigAs(mockFile)
	if err != nil {
		return err
	}

	return nil
}

//
// Recording Configuration Functions.
//

// RecordRequestHeaders define the request headers that should be recorded.
func RecordRequestHeaders(h ...string) RecordConfigurer {
	return recordConfigFunc(func(c *RecordConfig) { c.RequestHeaders = h })
}

// RecordResponseHeaders define the response headers that should be recorded.
func RecordResponseHeaders(h ...string) RecordConfigurer {
	return recordConfigFunc(func(c *RecordConfig) { c.ResponseHeaders = h })
}

// RecordDir defines the directory to save the recorded mocks.
// Defaults to: testdata/mocks
func RecordDir(dir string) RecordConfigurer {
	return recordConfigFunc(func(c *RecordConfig) { c.SaveDir = dir })
}

// RecordExtension defines the extension to save the recorded mock.
// It uses viper.Viper to save the mock files, so it will accept the same values that viper supports.
// Defaults to .json.
func RecordExtension(ext string) RecordConfigurer {
	return recordConfigFunc(func(c *RecordConfig) { c.SaveExtension = ext })
}

// RecordResponseBodyToFile defines if the recorded body should be saved to separate file.
// Defaults to false, embedding the response body into the mock definition file.
func RecordResponseBodyToFile(v bool) RecordConfigurer {
	return recordConfigFunc(func(c *RecordConfig) { c.SaveResponseBodyToFile = v })
}
