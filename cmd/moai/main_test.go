package main

import (
	"context"
	"net"
	"net/http"
	"os"
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/vitorsalgado/mocha/v3"
	"github.com/vitorsalgado/mocha/v3/internal/testutil"
)

func TestCLI(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	t.Cleanup(func() { cancel() })

	sc := mocha.StatusNoMatch

	port, err := getRandomPort()
	require.NoError(t, err)

	portTxt := strconv.FormatInt(int64(port), 10)
	cfg := &mocha.Config{Addr: ":" + portTxt}

	go run(ctx, cfg)

	request := func() {
		max := 3
		attempts := 0
		success := false

		for i := 0; i < max; i++ {
			res, err := testutil.Get("http://localhost:" + portTxt).Do()
			if err != nil || res.StatusCode == sc {
				success = true
				break
			}

			if attempts == max {
				t.FailNow()
			}

			<-time.After(2 * time.Second)

			attempts++
		}

		if !success {
			t.FailNow()
		}

	}

	request()

	cancel()
}

func TestDockerCLI(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	t.Cleanup(func() {
		cancel()
		_ = os.Unsetenv(_dockerHostEnv)
	})

	sc := http.StatusTeapot

	require.NoError(t, os.Setenv(_dockerHostEnv, "localhost"))

	go run(ctx)

	request := func() {
		max := 5
		attempts := 0
		success := false

		for i := 0; i < max; i++ {
			res, err := testutil.Get("http://localhost:8080").Do()

			if err == nil && res.StatusCode == sc {
				success = true
				break
			}

			if attempts == max {
				t.FailNow()
			}

			<-time.After(2 * time.Second)

			attempts++
		}

		if !success {
			t.FailNow()
		}

	}

	request()

	cancel()
}

func getRandomPort() (int, error) {
	l, err := net.Listen("tcp", ":0")
	if err != nil {
		return 0, err
	}

	port := l.Addr().(*net.TCPAddr).Port
	_ = l.Close()

	return port, nil
}
