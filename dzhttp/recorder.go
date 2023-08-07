package dzhttp

import (
	"bytes"
	"compress/gzip"
	"encoding/hex"
	"fmt"
	"hash/fnv"
	"io"
	"mime"
	"net/http"
	"net/url"
	"os"
	"path"
	"regexp"
	"strings"
	"sync"

	"github.com/spf13/viper"

	"github.com/vitorsalgado/mocha/v3/dzhttp/httpval"
)

var _ RecordConfigurer = (*RecordConfig)(nil)

var _regexNoSpecialCharacters = regexp.MustCompile("[^a-z0-9]")

type recorder struct {
	app    *HTTPMockApp
	active bool
	config *RecordConfig
	in     chan recArgs
	close  chan struct{}
	mu     sync.Mutex
}

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
	SaveDir string

	// SaveFileType defines the content type and extension to save the recorded mock.
	// It uses viper.Viper to save the mock files, so it will accept the same values that viper supports.
	// Defaults to json.
	SaveFileType string

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
	extension := strings.TrimPrefix(opts.SaveFileType, ".")
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

	opts.SaveFileType = extension
	opts.SaveDir = r.SaveDir
	opts.SaveResponseBodyToFile = r.SaveResponseBodyToFile
	opts.ResponseHeaders = r.ResponseHeaders
	opts.RequestHeaders = r.RequestHeaders

	return nil
}

func defaultRecordConfig() *RecordConfig {
	return &RecordConfig{
		SaveDir:                "testdata/_mocks_recorded",
		SaveFileType:           "json",
		SaveResponseBodyToFile: false,
		RequestHeaders:         []string{"accept", "content-type", "content-encoding", "content-length"},
		ResponseHeaders:        []string{"content-type", "content-encoding", "link", "content-length", "cache-control", "retry-after"},
	}
}

// Recording data transfer structs.
type recArgs struct {
	request  recRequest
	response recResponse
}

type recRequest struct {
	uri    string
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

func newRecorder(app *HTTPMockApp, config *RecordConfig) *recorder {
	return &recorder{app: app, config: config, close: make(chan struct{})}
}

func (r *recorder) start() {
	r.in = make(chan recArgs)
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
					r.app.logger.Error().Err(err).Msgf(err.Error())
				}

			case <-r.close:
				r.stop()
				return
			}
		}
	}()
}

func (r *recorder) stop() {
	r.mu.Lock()
	defer r.mu.Unlock()

	if !r.active {
		return
	}

	r.close <- struct{}{}
	r.active = false
}

func (r *recorder) dispatch(req *http.Request, parsedURL *url.URL, rawReqBody []byte, res *MockedResponse) {
	if !r.active {
		return
	}

	r.in <- recArgs{
		request:  recRequest{parsedURL.RequestURI(), req.Method, req.Header.Clone(), req.URL.Query(), rawReqBody},
		response: recResponse{res.StatusCode, res.Header.Clone(), res.Body},
	}
}

func (r *recorder) process(arg recArgs) error {
	v := viper.New()

	v.Set(_fRequestURLPath, arg.request.uri)
	v.Set(_fRequestMethod, arg.request.method)

	requestHeaders := make(map[string]string, len(arg.request.header))
	responseHeaders := make(map[string]string, len(arg.response.header))
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

	contentType := arg.request.header.Get(httpval.HeaderContentType)
	if strings.Contains(contentType, ";") {
		contentType = strings.TrimSpace(contentType[:strings.Index(contentType, ";")])
	}

	nm := strings.ReplaceAll(arg.request.uri, " ", "-")
	nm = strings.TrimPrefix(_regexNoSpecialCharacters.ReplaceAllString(nm, "-"), "-")
	if requestBodyHash != "" {
		nm += "--" + requestBodyHash
	}

	saveDir := r.config.SaveDir
	if !path.IsAbs(saveDir) {
		saveDir = path.Join(r.app.config.RootDir, saveDir)
	}

	name := fmt.Sprintf("%s-%s", arg.request.method, nm)
	mockFile := path.Join(saveDir, fmt.Sprintf("%s.mock.%s", name, r.config.SaveFileType))

	if hasResBody {
		if r.config.SaveResponseBodyToFile {
			bodyFile := name + "--response-body"
			ext, _ := mime.ExtensionsByType(contentType)

			if len(ext) > 0 {
				bodyFile += ext[0]
			} else {
				bodyFile += ".bin"
			}

			bodyFilename := path.Join(saveDir, bodyFile)
			b, err := os.Create(bodyFilename)
			if err != nil {
				return err
			}

			defer b.Close()

			// encoding := arg.response.header.Get(httpval.HeaderContentEncoding)
			// switch encoding {
			// case "gzip":
			// 	gz, err := gzip.NewReader(bytes.NewReader(arg.response.body))
			// 	if err != nil {
			// 		return err
			// 	}

			// 	defer gz.Close()

			// 	_, err = io.Copy(b, gz)
			// 	if err != nil {
			// 		return err
			// 	}

			// 	v.Set(_fResponseEncoding, "gzip")
			// default:
			// 	b.Write(arg.response.body)
			// }

			b.Write(arg.response.body)

			v.Set(_fResponseBodyFile, bodyFilename)

		} else {
			encoding := arg.response.header.Get(httpval.HeaderContentEncoding)
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

	v.Set(_fResponseStatus, arg.response.status)
	v.Set(_fResponseHeader, responseHeaders)

	v.AddConfigPath(saveDir)
	v.SetConfigType(r.config.SaveFileType)

	return v.WriteConfigAs(mockFile)
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
func RecordDir(dir string) RecordConfigurer {
	return recordConfigFunc(func(c *RecordConfig) { c.SaveDir = dir })
}

// RecordExtension defines the extension to save the recorded mock.
// It uses viper.Viper to save the mock files, so it will accept the same values that viper supports.
// Defaults to .json.
func RecordExtension(ext string) RecordConfigurer {
	return recordConfigFunc(func(c *RecordConfig) { c.SaveFileType = ext })
}

// RecordResponseBodyToFile defines if the recorded body should be saved to separate file.
// Defaults to false, embedding the response body into the mock definition file.
func RecordResponseBodyToFile(v bool) RecordConfigurer {
	return recordConfigFunc(func(c *RecordConfig) { c.SaveResponseBodyToFile = v })
}
