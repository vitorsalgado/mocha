//go:build docker
// +build docker

package testdocker

import (
	"crypto/tls"
	"crypto/x509"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

const (
	_dockerName        = "httpm_test"
	_tlsCertFile       = "../testdata/cert/cert.pem"
	_tlsKeyFile        = "../testdata/cert/key.pem"
	_tlsClientCertFile = "../testdata/cert/cert_client.pem"
)

var (
	// port, _    = testutil.RandomTCPPort()
	baseURL    = "http://localhost:8080"
	httpClient *http.Client
)

func TestMain(m *testing.M) {
	cert, err := tls.LoadX509KeyPair(_tlsCertFile, _tlsKeyFile)
	if err != nil {
		log.Fatalf("docker_test: failed to load tls certificate key pair: %v", err)
	}

	caCert, err := os.ReadFile(filepath.Clean(_tlsClientCertFile))
	if err != nil {
		log.Fatalf("docker_test: failed to load tls client certificate: %v", err)
	}

	certPool := x509.NewCertPool()
	certPool.AppendCertsFromPEM(caCert)

	httpClient = &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				Certificates: []tls.Certificate{cert},
				RootCAs:      certPool,
				ClientAuth:   tls.RequireAndVerifyClientCert,
			}},
	}

	// docker commands
	build := exec.Command("docker", "build", "-f", "../../Dockerfile", "-t", _dockerName, "../..")
	build.Stdout = os.Stdout
	build.Stderr = os.Stderr

	stop := exec.Command("docker", "stop", _dockerName)
	stop.Stdout = os.Stdout
	stop.Stderr = os.Stderr

	rm := exec.Command("docker", "rm", "-f", _dockerName)
	rm.Stdout = os.Stdout
	rm.Stderr = os.Stderr

	run := exec.Command(
		"docker",
		"run",
		"-d",
		"-it",
		// "-e MOAI_HTTPS=true",
		// "-e MOAI_TLS_CERT="+_tlsCertFile,
		// "-e MOAI_TLS_KEY="+_tlsKeyFile,
		// "-e MOAI_TLS_CA="+_tlsClientCertFile,
		"--network",
		"host",
		"--name",
		_dockerName,
		_dockerName,
	)
	run.Stdout = os.Stdout
	run.Stderr = os.Stderr

	// preparing test environment
	err = build.Run()
	if err != nil {
		log.Fatalf("test_docker: failed to build image %s: %v", _dockerName, err)
	}

	err = run.Run()
	if err != nil {
		log.Fatalf("test_docker: failed to run container %s: %v", _dockerName, err)
	}

	// running tests
	time.Sleep(10 * time.Second)
	code := m.Run()

	// cleanup
	err = stop.Run()
	if err != nil {
		log.Printf("test_docker: warning: failed to stop container %s: %v", _dockerName, err)
	}

	err = rm.Run()
	if err != nil {
		log.Printf("test_docker: warning: failed to remove container %s: %v", _dockerName, err)
	}

	// result
	os.Exit(code)
}

func TestBasic(t *testing.T) {
	res, err := httpClient.Get(buildURL(t, "/hello"))
	require.NoError(t, err)

	b, err := io.ReadAll(res.Body)

	require.NoError(t, err)
	require.NoError(t, res.Body.Close())
	require.Equal(t, http.StatusOK, res.StatusCode)
	require.Equal(t, "world", string(b))
}

// func TestCORS(t *testing.T) {
// 	corsReq, _ := http.NewRequest(http.MethodOptions, buildURL(t)+"/cors", nil)
// 	res, err := httpClient.Do(corsReq)
//
// 	require.NoError(t, err)
// 	require.Equal(t, http.StatusNoContent, res.StatusCode)
//
// 	req, _ := http.NewRequest(http.MethodGet, m.URL()+"/test", nil)
// 	res, err = httpClient.Do(req)
//
// 	require.NoError(t, err)
// 	require.Equal(t, http.StatusOK, res.StatusCode)
// }

func buildURL(t *testing.T, paths ...string) string {
	if len(paths) == 0 {
		return baseURL
	}

	u, err := url.JoinPath(baseURL, paths...)
	if err != nil {
		t.Errorf("server: building server url with path elements %s", paths)
	}

	return u
}
