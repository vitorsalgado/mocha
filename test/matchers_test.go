package test

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/vitorsalgado/mocha/v3"
)

func TestMatcherCombinations(t *testing.T) {
	type addr struct {
		Str     string  `json:"str,omitempty"`
		Number  int     `json:"number,omitempty"`
		Country *string `json:"country"`
	}

	type model struct {
		Name       string   `json:"name,omitempty"`
		Hobbies    []string `json:"hobbies,omitempty"`
		Job        string   `json:"job,omitempty"`
		Active     bool     `json:"active,omitempty"`
		Plans      []int    `json:"plans,omitempty"`
		Acronym    string   `json:"acronym,omitempty"`
		Title      string   `json:"title,omitempty"`
		Addr       addr     `json:"addr,omitempty"`
		CheckIn    bool     `json:"check_in,omitempty"`
		Activities string   `json:"activities,omitempty"`
		Department string   `json:"department,omitempty"`
		Place      string   `json:"place,omitempty"`
		Numbers    []int    `json:"numbers,omitempty"`

		Days int `json:"days,omitempty"`
	}

	type response struct {
		Ok   bool    `json:"ok,omitempty"`
		Type string  `json:"type,omitempty"`
		Num  float64 `json:"num,omitempty"`
	}

	httpClient := &http.Client{Transport: &http.Transport{TLSClientConfig: &tls.Config{InsecureSkipVerify: true}}}
	m := mocha.New()

	sHTTP := m.MustMock(mocha.FromFile("testdata/matchers/01_matchers.yaml"))
	sHTTPs := m.MustMock(mocha.FromFile("testdata/matchers/02_matchers.yaml"))

	actuateAndAssert := func(baseURL string) {
		body := &model{
			Name:    "hello world",
			Hobbies: []string{"bike", "trekking"},
			Job:     "dev",
			Active:  true,
			Acronym: "QA",
			Title:   "Software Engineer",
			Addr: addr{
				Str:     "Nowhere",
				Number:  100,
				Country: nil,
			},
			CheckIn:    false,
			Activities: "walk+run+swim",
			Department: "  Bombs  ",
			Place:      "Berlin",
			Numbers:    []int{10, 20, 30, 40, 50},

			Days: 10,
		}

		buf := new(bytes.Buffer)
		err := json.NewEncoder(buf).Encode(body)
		require.NoError(t, err)

		qry := url.Values{}
		qry.Add("q", "test")
		qry.Add("sort", "asc")
		qry.Add("page", "5")
		qry.Add("lang", "en")
		qry.Add("test", "true")
		qry.Add("tags", "link,detail,price")
		qry.Add("ctx", "hello world")
		qry.Add("nm", "hello")
		qry.Add("nm", "hi")
		qry.Add("nm", "hallo")

		u := baseURL + "/test?" + qry.Encode() + "&age(150)&contains(nothing)"
		req, _ := http.NewRequest(http.MethodPost, u, buf)
		req.Header.Add(mocha.HeaderContentType, mocha.MIMEApplicationJSON)
		req.Header.Add(mocha.HeaderAccept, mocha.MIMEApplicationJSON)
		req.Header.Add("x-test", "dev, qa, devops")

		res, err := httpClient.Do(req)
		require.NoError(t, err)

		b, err := io.ReadAll(res.Body)
		require.NoError(t, err)
		require.NoError(t, res.Body.Close())

		r := new(response)
		err = json.Unmarshal(b, r)

		require.NoError(t, err, baseURL)
		require.Equal(t, mocha.MIMEApplicationJSON, res.Header.Get(mocha.HeaderContentType))
		require.Equal(t, "success", res.Header.Get("x-custom"))
		require.Equal(t, http.StatusOK, res.StatusCode)
		require.Equal(t, true, r.Ok)
		require.Equal(t, "test", r.Type)
		require.EqualValues(t, 10, r.Num)
	}

	m.MustStart()
	u := m.URL()
	actuateAndAssert(u)
	m.Close()

	m.MustStartTLS()
	u = m.URL()
	actuateAndAssert(u)
	m.Close()

	require.True(t, sHTTP.AssertNumberOfCalls(t, 1))
	require.True(t, sHTTPs.AssertNumberOfCalls(t, 1))
}

func TestMatchers_MultipleMethods(t *testing.T) {
	m := mocha.NewT(t)
	m.MustStart()
	m.MustMock(mocha.FromFile("testdata/matchers/03_url_template.yaml"))

	client := &http.Client{}
	res, err := client.Get(m.URL("/test?q=none"))

	require.NoError(t, err)
	require.Equal(t, http.StatusNoContent, res.StatusCode)

	res, err = client.Post(m.URL("/test?q=none"), mocha.MIMETextPlain, nil)

	require.NoError(t, err)
	require.Equal(t, http.StatusNoContent, res.StatusCode)

	req, _ := http.NewRequest(http.MethodPut, m.URL("/test?q=none"), nil)
	res, err = client.Do(req)

	require.NoError(t, err)
	require.Equal(t, mocha.StatusNoMatch, res.StatusCode)
}

func TestMatchers_URLMatchRegex(t *testing.T) {
	m := mocha.NewT(t)
	m.MustStart()
	m.MustMock(mocha.FromFile("testdata/matchers/04_url_match.yaml"))

	client := &http.Client{}
	res, err := client.Get(m.URL("/test?q=hi"))

	require.NoError(t, err)
	require.Equal(t, http.StatusNoContent, res.StatusCode)

	res, err = client.Get(m.URL("/test?q=bye"))

	require.NoError(t, err)
	require.Equal(t, mocha.StatusNoMatch, res.StatusCode)
}

func TestMatchers_URLCustomMatcher(t *testing.T) {
	m := mocha.NewT(t)
	m.MustStart()
	m.MustMock(mocha.FromFile("testdata/matchers/05_url_custom_matcher.yaml"))

	client := &http.Client{}
	res, err := client.Get(m.URL("/test?q=hi"))

	require.NoError(t, err)
	require.Equal(t, http.StatusNoContent, res.StatusCode)

	res, err = client.Get(m.URL("/test?q=bye"))

	require.NoError(t, err)
	require.Equal(t, mocha.StatusNoMatch, res.StatusCode)
}

func TestMatchers_PathCustomMatcher(t *testing.T) {
	m := mocha.NewT(t)
	m.MustStart()
	m.MustMock(mocha.FromFile("testdata/matchers/06_path_custom_matcher.yaml"))

	client := &http.Client{}
	res, err := client.Get(m.URL("/hi"))

	require.NoError(t, err)
	require.Equal(t, http.StatusNoContent, res.StatusCode)

	res, err = client.Get(m.URL("/bye"))

	require.NoError(t, err)
	require.Equal(t, mocha.StatusNoMatch, res.StatusCode)
}

func TestMatchers_NoReply_ShouldReturn200ByDefault(t *testing.T) {
	m := mocha.NewT(t)
	m.MustStart()
	m.MustMock(mocha.FromFile("testdata/matchers/07_no_reply.yaml"))

	client := &http.Client{}
	res, err := client.Get(m.URL("/test"))

	require.NoError(t, err)
	require.Equal(t, http.StatusOK, res.StatusCode)
}
