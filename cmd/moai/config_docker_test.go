package main

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/vitorsalgado/mocha/v3"
)

func TestConfigDocker_DefaultHTTPHost(t *testing.T) {
	config := &mocha.Config{}
	conf := &dockerConfigurer{}
	err := conf.Apply(config)

	assert.NoError(t, err)
	assert.False(t, config.UseHTTPS)
	assert.Equal(t, "0.0.0.0:8080", config.Addr)
}

func TestConfigDocker_DefaultHTTPsHost(t *testing.T) {
	config := &mocha.Config{}
	config.UseHTTPS = true

	conf := &dockerConfigurer{}
	err := conf.Apply(config)

	assert.NoError(t, err)
	assert.True(t, config.UseHTTPS)
	assert.Equal(t, "0.0.0.0:8443", config.Addr)
}

func TestConfigDocker_HostDefined(t *testing.T) {
	t.Cleanup(func() {
		_ = os.Unsetenv(_dockerHostEnv)
	})

	require.NoError(t, os.Setenv(_dockerHostEnv, "example.org"))

	config := &mocha.Config{}
	conf := &dockerConfigurer{}
	err := conf.Apply(config)

	assert.NoError(t, err)
	assert.False(t, config.UseHTTPS)
	assert.Equal(t, "example.org:8080", config.Addr)
}

func TestConfigDocker_HostDefined_WithDots(t *testing.T) {
	t.Cleanup(func() {
		_ = os.Unsetenv(_dockerHostEnv)
	})

	require.NoError(t, os.Setenv(_dockerHostEnv, "example.org:"))

	config := &mocha.Config{}
	conf := &dockerConfigurer{}
	err := conf.Apply(config)

	assert.NoError(t, err)
	assert.False(t, config.UseHTTPS)
	assert.Equal(t, "example.org:8080", config.Addr)
}
