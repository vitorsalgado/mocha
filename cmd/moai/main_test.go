package main

import (
	"context"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/vitorsalgado/mocha/v3/internal/testutil"
)

func TestCLI(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	t.Cleanup(func() { cancel() })

	sc := http.StatusTeapot

	go run(ctx)

	request := func() {
		max := 3
		attempts := 0
		success := false

		for i := 0; i < max; i++ {
			res, err := testutil.Get("http://localhost:3000").Do()
			require.NoError(t, err)

			if res.StatusCode == sc {
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
