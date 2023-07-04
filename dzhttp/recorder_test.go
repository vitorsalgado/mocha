package dzhttp

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"os"
	"path"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	. "github.com/vitorsalgado/mocha/v3/matcher"
)

func TestRecordConfigApply(t *testing.T) {
	conf := &RecordConfig{}

	require.Error(t, conf.Apply(&RecordConfig{}), "no extension defined")
	require.Error(t, conf.Apply(&RecordConfig{SaveFileType: ".test"}), "unsupported extension")
	require.NoError(t, conf.Apply(&RecordConfig{SaveFileType: ".json"}), "should accept extension with dot prefix")
	require.NoError(t, conf.Apply(&RecordConfig{SaveFileType: "json"}), "should accept extension without dot prefix")
}

func TestRecordingWithWebProxy(t *testing.T) {
	dir := t.TempDir()

	includedReqHeader := "x-request-excluded"
	excludedReqHeader := "x-req-excluded"

	includedResHeader := "x-response-included"
	excludedResHeader := "x-res-excluded"

	trailer := "final"

	p := NewAPI(Setup().
		Name("recorder").
		Proxy().
		Record(
			RecordDir(dir),
			RecordResponseBodyToFile(true),
			RecordRequestHeaders("accept", "content-type", "content-encoding", "content-length", includedReqHeader),
			RecordResponseHeaders("content-type", "content-encoding", "link", "content-length", includedResHeader, trailer),
		))
	p.MustStart()
	scope1 := p.MustMock(Get(URLPath("/test")).Reply(Accepted()))

	m := NewAPI()
	m.MustStart()
	scope2 := m.MustMock(Get(URLPath("/other")).
		Reply(Created().
			Header(includedResHeader, "included").
			Header(excludedResHeader, "excluded").
			Trailer(trailer, "the-trailer-value").
			BodyText("hello world")),
	)

	u, _ := url.Parse(p.URL())
	client := &http.Client{Transport: &http.Transport{Proxy: http.ProxyURL(u)}}

	req, _ := http.NewRequest(http.MethodGet, m.URL()+"/test", nil)
	res, err := client.Do(req)

	require.NoError(t, err)
	require.NoError(t, res.Body.Close())
	require.Equal(t, http.StatusAccepted, res.StatusCode)

	req, _ = http.NewRequest(http.MethodGet, m.URL()+"/other", nil)
	req.Header.Add(excludedReqHeader, "excluded")
	req.Header.Add(includedReqHeader, "included")
	res, err = client.Do(req)

	require.NoError(t, err)
	require.NoError(t, res.Body.Close())
	require.Equal(t, http.StatusCreated, res.StatusCode)

	scope1.AssertCalled(t)
	scope2.AssertCalled(t)

	time.Sleep(2 * time.Second)

	entries, err := os.ReadDir(dir)

	require.NoError(t, err)
	require.Len(t, entries, 2)

	p.Close()
	m.Close()

	// Creating a new server that will use the recorded mocks

	srv := NewAPI(Setup().MockFilePatterns(dir + "/*mock.json"))
	srv.MustStart()

	defer srv.Close()

	httpClient := &http.Client{}

	req, _ = http.NewRequest(http.MethodGet, srv.URL()+"/other", nil)
	req.Header.Add(excludedReqHeader, "excluded")
	res, err = httpClient.Do(req)

	require.NoError(t, err)
	require.Equal(t, StatusNoMatch, res.StatusCode)

	req, _ = http.NewRequest(http.MethodGet, srv.URL()+"/other", nil)
	req.Header.Add(includedReqHeader, "included")
	res, err = httpClient.Do(req)
	require.NoError(t, err)

	b, err := io.ReadAll(res.Body)
	res.Body.Close()

	require.NoError(t, err)
	require.Equal(t, http.StatusCreated, res.StatusCode)
	require.Equal(t, "hello world", string(b))
	require.Empty(t, res.Header.Get(excludedResHeader))
	require.Equal(t, res.Header.Get(includedResHeader), "included")
	require.Equal(t, res.Header.Get(trailer), "the-trailer-value")
}

func TestRecordingWithWebProxy_CustomRootDir(t *testing.T) {
	dir := t.TempDir()
	actualDir := path.Join(dir, "testdata/_mocks_recorded")

	err := os.MkdirAll(actualDir, os.ModePerm)
	require.NoError(t, err)

	includedReqHeader := "x-request-excluded"
	excludedReqHeader := "x-req-excluded"

	includedResHeader := "x-response-included"
	excludedResHeader := "x-res-excluded"

	trailer := "final"

	p := NewAPI(Setup().
		Name("recorder").
		RootDir(dir).
		Proxy().
		Record(
			RecordResponseBodyToFile(true),
			RecordRequestHeaders("accept", "content-type", "content-encoding", "content-length", includedReqHeader),
			RecordResponseHeaders("content-type", "content-encoding", "link", "content-length", includedResHeader, trailer),
		))
	p.MustStart()
	scope1 := p.MustMock(Get(URLPath("/test")).Reply(Accepted()))

	m := NewAPI()
	m.MustStart()
	scope2 := m.MustMock(Get(URLPath("/other")).
		Reply(Created().
			Header(includedResHeader, "included").
			Header(excludedResHeader, "excluded").
			Trailer(trailer, "the-trailer-value").
			BodyText("hello world")),
	)

	u, _ := url.Parse(p.URL())
	client := &http.Client{Transport: &http.Transport{Proxy: http.ProxyURL(u)}}

	req, _ := http.NewRequest(http.MethodGet, m.URL()+"/test", nil)
	res, err := client.Do(req)

	require.NoError(t, err)
	require.NoError(t, res.Body.Close())
	require.Equal(t, http.StatusAccepted, res.StatusCode)

	req, _ = http.NewRequest(http.MethodGet, m.URL()+"/other", nil)
	req.Header.Add(excludedReqHeader, "excluded")
	req.Header.Add(includedReqHeader, "included")
	res, err = client.Do(req)

	require.NoError(t, err)
	require.NoError(t, res.Body.Close())
	require.Equal(t, http.StatusCreated, res.StatusCode)

	scope1.AssertCalled(t)
	scope2.AssertCalled(t)

	time.Sleep(2 * time.Second)

	entries, err := os.ReadDir(actualDir)

	require.NoError(t, err)
	require.Len(t, entries, 2)

	p.Close()
	m.Close()

	// Creating a new server that will use the recorded mocks

	srv := NewAPI(Setup().MockFilePatterns(actualDir + "/*mock.json"))
	srv.MustStart()

	defer srv.Close()

	httpClient := &http.Client{}

	req, _ = http.NewRequest(http.MethodGet, srv.URL()+"/other", nil)
	req.Header.Add(excludedReqHeader, "excluded")
	res, err = httpClient.Do(req)

	require.NoError(t, err)
	require.Equal(t, StatusNoMatch, res.StatusCode)

	req, _ = http.NewRequest(http.MethodGet, srv.URL()+"/other", nil)
	req.Header.Add(includedReqHeader, "included")
	res, err = httpClient.Do(req)
	require.NoError(t, err)

	b, err := io.ReadAll(res.Body)
	res.Body.Close()

	require.NoError(t, err)
	require.Equal(t, http.StatusCreated, res.StatusCode)
	require.Equal(t, "hello world", string(b))
	require.Empty(t, res.Header.Get(excludedResHeader))
	require.Equal(t, res.Header.Get(includedResHeader), "included")
	require.Equal(t, res.Header.Get(trailer), "the-trailer-value")
}

func TestRecordSaveResponseBodyToFile(t *testing.T) {
	dir := t.TempDir()
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	file, err := os.Open(path.Join("testdata/recorder/customers.json"))
	require.NoError(t, err)

	defer file.Close()

	data := make([]map[string]any, 0)
	err = json.NewDecoder(file).Decode(&data)
	require.NoError(t, err)

	target := NewAPI(Setup().Name("target"))
	target.MustStart()
	targetScope := target.MustMock(
		Get(URLPath("/customers")).Reply(OK().JSON(data).Gzip()),
		Getf("/customers/001").Reply(OK().JSON(data[0]).Gzip()),
		Getf("/customers/002").Reply(OK().JSON(data[1])),
		Postf("/customers").Reply(Created()),
	)

	rec := NewAPI(Setup().Name("recorder").Record(
		RecordDir(dir),
		RecordResponseBodyToFile(true),
	))
	rec.MustStart()
	recorderScope := rec.MustMock(AnyMethod().Reply(From(target.URL())))

	httpClient := &http.Client{}

	res, err := httpClient.Get(rec.URL() + "/customers")
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, res.StatusCode)

	b := make([]map[string]any, 0)
	err = json.NewDecoder(res.Body).Decode(&b)
	require.NoError(t, err)
	require.Equal(t, data, b)

	res.Body.Close()

	res, err = httpClient.Get(rec.URL() + "/customers/001")
	require.NoError(t, err)
	require.NoError(t, res.Body.Close())
	require.Equal(t, http.StatusOK, res.StatusCode)

	res, err = httpClient.Get(rec.URL() + "/customers/002")
	require.NoError(t, err)

	b2 := make(map[string]any)
	err = json.NewDecoder(res.Body).Decode(&b2)
	require.NoError(t, err)

	res.Body.Close()

	require.Equal(t, http.StatusOK, res.StatusCode)
	require.Equal(t, data[1], b2)

	res, err = httpClient.Post(rec.URL()+"/customers", "application/json", nil)
	require.NoError(t, err)
	require.NoError(t, res.Body.Close())
	require.Equal(t, http.StatusCreated, res.StatusCode)

	recorderScope.AssertCalled(t)
	targetScope.AssertCalled(t)

	<-ctx.Done()

	entries, err := os.ReadDir(dir)

	require.NoError(t, err)
	require.Len(t, entries, 7)

	rec.Close()
	target.Close()

	// Creating a new server that will use the recorded mocks

	m := NewAPI(Setup().MockFilePatterns(dir + "/*mock.json"))
	m.MustStart()

	defer m.Close()

	res, err = httpClient.Get(m.URL() + "/customers")
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, res.StatusCode)

	rb := make([]map[string]any, 0)
	err = json.NewDecoder(res.Body).Decode(&rb)
	require.NoError(t, err)
	assert.Equal(t, data, rb)

	res.Body.Close()

	res, err = httpClient.Get(m.URL() + "/customers/001")
	require.NoError(t, err)
	require.NoError(t, res.Body.Close())
	assert.Equal(t, http.StatusOK, res.StatusCode)

	res, err = httpClient.Get(m.URL() + "/customers/002")
	require.NoError(t, err)

	rb2 := make(map[string]any)
	err = json.NewDecoder(res.Body).Decode(&rb2)
	require.NoError(t, err)

	res.Body.Close()

	assert.Equal(t, http.StatusOK, res.StatusCode)
	assert.Equal(t, data[1], rb2)

	res, err = httpClient.Post(m.URL()+"/customers", "application/json", nil)
	require.NoError(t, err)
	require.NoError(t, res.Body.Close())
	require.Equal(t, http.StatusCreated, res.StatusCode)
}

func TestRecordEmbeddedResponseBodiesYAML(t *testing.T) {
	dir := t.TempDir()
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	file, err := os.Open(path.Join("testdata/recorder/customers.json"))
	require.NoError(t, err)

	defer file.Close()

	data := make([]map[string]any, 0)
	err = json.NewDecoder(file).Decode(&data)
	require.NoError(t, err)

	target := NewAPI(Setup().Name("target"))
	target.MustStart()
	targetScope := target.MustMock(
		Get(URLPath("/customers")).Reply(OK().JSON(data).Gzip()),
		Getf("/customers/001").Reply(OK().JSON(data[0]).Gzip()),
		Getf("/customers/002").Reply(OK().JSON(data[1])),
		Postf("/customers").Reply(Created()),
	)

	rec := NewAPI(Setup().Name("recorder").Record(
		RecordDir(dir),
		RecordResponseBodyToFile(false),
		RecordExtension("yaml"),
	))
	rec.MustStart()
	recorderScope := rec.MustMock(AnyMethod().Reply(From(target.URL())))

	httpClient := &http.Client{}

	res, err := httpClient.Get(rec.URL() + "/customers")
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, res.StatusCode)

	b := make([]map[string]any, 0)
	err = json.NewDecoder(res.Body).Decode(&b)
	require.NoError(t, err)
	require.Equal(t, data, b)

	res.Body.Close()

	res, err = httpClient.Get(rec.URL() + "/customers/001")
	require.NoError(t, err)
	require.NoError(t, res.Body.Close())
	require.Equal(t, http.StatusOK, res.StatusCode)

	res, err = httpClient.Get(rec.URL() + "/customers/002")
	require.NoError(t, err)

	defer res.Body.Close()

	b2 := make(map[string]any)
	err = json.NewDecoder(res.Body).Decode(&b2)
	require.NoError(t, err)

	assert.Equal(t, http.StatusOK, res.StatusCode)
	assert.Equal(t, data[1], b2)

	res, err = httpClient.Post(rec.URL()+"/customers", "application/json", nil)
	require.NoError(t, err)
	require.NoError(t, res.Body.Close())
	assert.Equal(t, http.StatusCreated, res.StatusCode)

	recorderScope.AssertCalled(t)
	targetScope.AssertCalled(t)

	<-ctx.Done()

	entries, err := os.ReadDir(dir)

	require.NoError(t, err)
	assert.Len(t, entries, 4)

	rec.Close()
	target.Close()

	// Creating a new server that will use the recorded mocks

	m := NewAPI(Setup().MockFilePatterns(dir + "/*mock.yaml"))
	m.MustStart()

	defer m.Close()

	res, err = httpClient.Get(m.URL() + "/customers")
	require.NoError(t, err)

	defer res.Body.Close()

	assert.Equal(t, http.StatusOK, res.StatusCode)

	rb := make([]map[string]any, 0)
	err = json.NewDecoder(res.Body).Decode(&rb)
	require.NoError(t, err)
	assert.Equal(t, data, rb)

	res, err = httpClient.Get(m.URL() + "/customers/001")
	require.NoError(t, err)
	require.NoError(t, res.Body.Close())
	assert.Equal(t, http.StatusOK, res.StatusCode)

	res, err = httpClient.Get(m.URL() + "/customers/002")
	require.NoError(t, err)

	defer res.Body.Close()

	rb2 := make(map[string]any)
	err = json.NewDecoder(res.Body).Decode(&rb2)
	require.NoError(t, err)

	assert.Equal(t, http.StatusOK, res.StatusCode)
	assert.Equal(t, data[1], rb2)

	res, err = httpClient.Post(m.URL()+"/customers", "application/json", nil)
	require.NoError(t, err)
	require.NoError(t, res.Body.Close())
	assert.Equal(t, http.StatusCreated, res.StatusCode)
}

func TestRecordTargetTLS(t *testing.T) {
	dir := t.TempDir()
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	file, err := os.Open(path.Join("testdata/recorder/customers.json"))
	require.NoError(t, err)

	defer file.Close()

	data := make([]map[string]any, 0)
	err = json.NewDecoder(file).Decode(&data)
	require.NoError(t, err)

	target := NewAPI(Setup().Name("target"))
	target.MustStartTLS()
	targetScope := target.MustMock(
		Get(URLPath("/customers")).Reply(OK().JSON(data).Gzip()),
		Getf("/customers/001").Reply(OK().JSON(data[0]).Gzip()),
		Getf("/customers/002").Reply(OK().JSON(data[1])),
		Postf("/customers").Reply(Created()),
	)

	recorder := NewAPI(Setup().Name("recorder").Record(
		RecordDir(dir),
		RecordResponseBodyToFile(false),
	))
	recorder.MustStart()
	recorderScope := recorder.MustMock(AnyMethod().Reply(From(target.URL()).SkipSSLVerify()))

	httpClient := &http.Client{}

	res, err := httpClient.Get(recorder.URL() + "/customers")
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, res.StatusCode)

	b := make([]map[string]any, 0)
	err = json.NewDecoder(res.Body).Decode(&b)
	require.NoError(t, err)
	require.Equal(t, data, b)

	res.Body.Close()

	res, err = httpClient.Get(recorder.URL() + "/customers/001")
	require.NoError(t, err)
	require.NoError(t, res.Body.Close())
	require.Equal(t, http.StatusOK, res.StatusCode)

	res, err = httpClient.Get(recorder.URL() + "/customers/002")
	require.NoError(t, err)

	defer res.Body.Close()

	b2 := make(map[string]any)
	err = json.NewDecoder(res.Body).Decode(&b2)
	require.NoError(t, err)

	assert.Equal(t, http.StatusOK, res.StatusCode)
	assert.Equal(t, data[1], b2)

	res, err = httpClient.Post(recorder.URL()+"/customers", "application/json", nil)
	require.NoError(t, err)
	require.NoError(t, res.Body.Close())
	assert.Equal(t, http.StatusCreated, res.StatusCode)

	recorderScope.AssertCalled(t)
	targetScope.AssertCalled(t)

	<-ctx.Done()

	entries, err := os.ReadDir(dir)

	require.NoError(t, err)
	assert.Len(t, entries, 4)

	recorder.Close()
	target.Close()

	// Creating a new server that will use the recorded mocks

	m := NewAPI(Setup().MockFilePatterns(dir + "/*mock.json"))
	m.MustStart()

	defer m.Close()

	res, err = httpClient.Get(m.URL() + "/customers")
	require.NoError(t, err)

	defer res.Body.Close()

	assert.Equal(t, http.StatusOK, res.StatusCode)

	rb := make([]map[string]any, 0)
	err = json.NewDecoder(res.Body).Decode(&rb)
	require.NoError(t, err)
	assert.Equal(t, data, rb)

	res, err = httpClient.Get(m.URL() + "/customers/001")
	require.NoError(t, err)
	require.NoError(t, res.Body.Close())
	assert.Equal(t, http.StatusOK, res.StatusCode)

	res, err = httpClient.Get(m.URL() + "/customers/002")
	require.NoError(t, err)

	defer res.Body.Close()

	rb2 := make(map[string]any)
	err = json.NewDecoder(res.Body).Decode(&rb2)
	require.NoError(t, err)

	assert.Equal(t, http.StatusOK, res.StatusCode)
	assert.Equal(t, data[1], rb2)

	res, err = httpClient.Post(m.URL()+"/customers", "application/json", nil)
	require.NoError(t, err)
	require.NoError(t, res.Body.Close())
	assert.Equal(t, http.StatusCreated, res.StatusCode)
}

func TestRecordTLSClient(t *testing.T) {
	dir := t.TempDir()
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	file, err := os.Open(path.Join("testdata/recorder/customers.json"))
	require.NoError(t, err)

	defer file.Close()

	data := make([]map[string]any, 0)
	err = json.NewDecoder(file).Decode(&data)
	require.NoError(t, err)

	target := NewAPI(Setup().Name("target"))
	target.MustStart()
	targetScope := target.MustMock(
		Get(URLPath("/customers")).Reply(OK().JSON(data).Gzip()),
		Getf("/customers/001").Reply(OK().JSON(data[0]).Gzip()),
		Getf("/customers/002").Reply(OK().JSON(data[1])),
		Postf("/customers").Reply(Created()),
	)

	recorder := NewAPI(Setup().Name("recorder").Record(
		RecordDir(dir),
		RecordResponseBodyToFile(false),
	))
	recorder.MustStartTLS()
	recorderScope := recorder.MustMock(AnyMethod().Reply(From(target.URL())))

	httpClient := &http.Client{Transport: &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}}

	res, err := httpClient.Get(recorder.URL() + "/customers")
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, res.StatusCode)

	b := make([]map[string]any, 0)
	err = json.NewDecoder(res.Body).Decode(&b)
	require.NoError(t, err)
	require.Equal(t, data, b)

	res.Body.Close()

	res, err = httpClient.Get(recorder.URL() + "/customers/001")
	require.NoError(t, err)
	require.NoError(t, res.Body.Close())
	require.Equal(t, http.StatusOK, res.StatusCode)

	res, err = httpClient.Get(recorder.URL() + "/customers/002")
	require.NoError(t, err)

	defer res.Body.Close()

	b2 := make(map[string]any)
	err = json.NewDecoder(res.Body).Decode(&b2)
	require.NoError(t, err)

	assert.Equal(t, http.StatusOK, res.StatusCode)
	assert.Equal(t, data[1], b2)

	res, err = httpClient.Post(recorder.URL()+"/customers", "application/json", nil)
	require.NoError(t, err)
	require.NoError(t, res.Body.Close())
	assert.Equal(t, http.StatusCreated, res.StatusCode)

	recorderScope.AssertCalled(t)
	targetScope.AssertCalled(t)

	<-ctx.Done()

	entries, err := os.ReadDir(dir)

	require.NoError(t, err)
	assert.Len(t, entries, 4)

	recorder.Close()
	target.Close()

	// Creating a new server that will use the recorded mocks

	m := NewAPI(Setup().MockFilePatterns(dir + "/*mock.json"))
	m.MustStart()

	defer m.Close()

	res, err = httpClient.Get(m.URL() + "/customers")
	require.NoError(t, err)

	defer res.Body.Close()

	assert.Equal(t, http.StatusOK, res.StatusCode)

	rb := make([]map[string]any, 0)
	err = json.NewDecoder(res.Body).Decode(&rb)
	require.NoError(t, err)
	assert.Equal(t, data, rb)

	res, err = httpClient.Get(m.URL() + "/customers/001")
	require.NoError(t, err)
	require.NoError(t, res.Body.Close())
	assert.Equal(t, http.StatusOK, res.StatusCode)

	res, err = httpClient.Get(m.URL() + "/customers/002")
	require.NoError(t, err)

	defer res.Body.Close()

	rb2 := make(map[string]any)
	err = json.NewDecoder(res.Body).Decode(&rb2)
	require.NoError(t, err)

	assert.Equal(t, http.StatusOK, res.StatusCode)
	assert.Equal(t, data[1], rb2)

	res, err = httpClient.Post(m.URL()+"/customers", "application/json", nil)
	require.NoError(t, err)
	require.NoError(t, res.Body.Close())
	assert.Equal(t, http.StatusCreated, res.StatusCode)
}

func TestRecordBothTLS(t *testing.T) {
	dir := t.TempDir()
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	file, err := os.Open(path.Join("testdata/recorder/customers.json"))
	require.NoError(t, err)

	defer file.Close()

	data := make([]map[string]any, 0)
	err = json.NewDecoder(file).Decode(&data)
	require.NoError(t, err)

	target := NewAPI(Setup().Name("target"))
	target.MustStartTLS()
	targetScope := target.MustMock(
		Get(URLPath("/customers")).Reply(OK().JSON(data).Gzip()),
		Getf("/customers/001").Reply(OK().JSON(data[0]).Gzip()),
		Getf("/customers/002").Reply(OK().JSON(data[1])),
		Postf("/customers").Reply(Created()),
	)

	recorder := NewAPI(Setup().Name("recorder").Record(
		RecordDir(dir),
		RecordResponseBodyToFile(false),
	))
	recorder.MustStartTLS()
	recorderScope := recorder.MustMock(AnyMethod().Reply(From(target.URL()).SkipSSLVerify()))

	httpClient := &http.Client{Transport: &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}}

	res, err := httpClient.Get(recorder.URL() + "/customers")
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, res.StatusCode)

	b := make([]map[string]any, 0)
	err = json.NewDecoder(res.Body).Decode(&b)
	require.NoError(t, err)
	require.Equal(t, data, b)

	res.Body.Close()

	res, err = httpClient.Get(recorder.URL() + "/customers/001")
	require.NoError(t, err)
	require.NoError(t, res.Body.Close())
	require.Equal(t, http.StatusOK, res.StatusCode)

	res, err = httpClient.Get(recorder.URL() + "/customers/002")
	require.NoError(t, err)

	defer res.Body.Close()

	b2 := make(map[string]any)
	err = json.NewDecoder(res.Body).Decode(&b2)
	require.NoError(t, err)

	assert.Equal(t, http.StatusOK, res.StatusCode)
	assert.Equal(t, data[1], b2)

	res, err = httpClient.Post(recorder.URL()+"/customers", "application/json", nil)
	require.NoError(t, err)
	require.NoError(t, res.Body.Close())
	assert.Equal(t, http.StatusCreated, res.StatusCode)

	recorderScope.AssertCalled(t)
	targetScope.AssertCalled(t)

	<-ctx.Done()

	entries, err := os.ReadDir(dir)

	require.NoError(t, err)
	assert.Len(t, entries, 4)

	recorder.Close()
	target.Close()

	// Creating a new server that will use the recorded mocks

	m := NewAPI(Setup().MockFilePatterns(dir + "/*mock.json"))
	m.MustStart()

	defer m.Close()

	res, err = httpClient.Get(m.URL() + "/customers")
	require.NoError(t, err)

	defer res.Body.Close()

	assert.Equal(t, http.StatusOK, res.StatusCode)

	rb := make([]map[string]any, 0)
	err = json.NewDecoder(res.Body).Decode(&rb)
	require.NoError(t, err)
	assert.Equal(t, data, rb)

	res, err = httpClient.Get(m.URL() + "/customers/001")
	require.NoError(t, err)
	require.NoError(t, res.Body.Close())
	assert.Equal(t, http.StatusOK, res.StatusCode)

	res, err = httpClient.Get(m.URL() + "/customers/002")
	require.NoError(t, err)

	defer res.Body.Close()

	rb2 := make(map[string]any)
	err = json.NewDecoder(res.Body).Decode(&rb2)
	require.NoError(t, err)

	assert.Equal(t, http.StatusOK, res.StatusCode)
	assert.Equal(t, data[1], rb2)

	res, err = httpClient.Post(m.URL()+"/customers", "application/json", nil)
	require.NoError(t, err)
	require.NoError(t, res.Body.Close())
	assert.Equal(t, http.StatusCreated, res.StatusCode)
}
