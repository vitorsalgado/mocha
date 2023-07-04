package main

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/vitorsalgado/mocha/v3/dzhttp"
)

func TestConfigDockerDefaultHTTPHost(t *testing.T) {
	config := &dzhttp.Config{}
	conf := &dockerConfigurer{}
	err := conf.Apply(config)

	assert.NoError(t, err)
	assert.False(t, config.UseHTTPS)
	assert.Equal(t, "0.0.0.0:8080", config.Addr)
}

func TestConfigDockerDefaultHTTPsHost(t *testing.T) {
	config := &dzhttp.Config{}
	config.UseHTTPS = true

	conf := &dockerConfigurer{}
	err := conf.Apply(config)

	assert.NoError(t, err)
	assert.True(t, config.UseHTTPS)
	assert.Equal(t, "0.0.0.0:8443", config.Addr)
}

func TestConfigDockerHostDefined(t *testing.T) {
	t.Cleanup(func() {
		_ = os.Unsetenv(_dockerHostEnv)
	})

	require.NoError(t, os.Setenv(_dockerHostEnv, "example.org"))

	config := &dzhttp.Config{}
	conf := &dockerConfigurer{}
	err := conf.Apply(config)

	assert.NoError(t, err)
	assert.False(t, config.UseHTTPS)
	assert.Equal(t, "example.org:8080", config.Addr)
}

func TestConfigDockerHostDefinedWithDots(t *testing.T) {
	t.Cleanup(func() {
		_ = os.Unsetenv(_dockerHostEnv)
	})

	require.NoError(t, os.Setenv(_dockerHostEnv, "example.org:"))

	config := &dzhttp.Config{}
	conf := &dockerConfigurer{}
	err := conf.Apply(config)

	assert.NoError(t, err)
	assert.False(t, config.UseHTTPS)
	assert.Equal(t, "example.org:8080", config.Addr)
}
