package mocha

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestViperConfigLoader_Load(t *testing.T) {
	cl := BuiltInConfigurer()
	config := &Config{}
	err := cl.Apply(config)
	require.NoError(t, err)

	assert.Equal(t, "test_api", config.Name)
	assert.Equal(t, ":8080", config.Addr)

	assert.NotNil(t, config.CORS)
	assert.Equal(t, "https://example.org", config.CORS.AllowedOrigin)
	assert.Equal(t, "GET,POST", config.CORS.AllowedMethods)
	assert.Equal(t, "*", config.CORS.AllowedHeaders)
	assert.Equal(t, "None", config.CORS.ExposeHeaders)
	assert.Equal(t, 150, config.CORS.MaxAge)
	assert.Equal(t, 200, config.CORS.SuccessStatusCode)

	assert.NotNil(t, config.Proxy)
	assert.Equal(t, "https://proxy.org/test", config.Proxy.Target)
	assert.Equal(t, time.Duration(5000), config.Proxy.Timeout)

	assert.NotNil(t, config.Record)
	assert.Equal(t, []string{"header1", "header2"}, config.Record.RequestHeaders)
	assert.Equal(t, []string{"header3", "header4"}, config.Record.ResponseHeaders)
	assert.True(t, config.Record.Save)
	assert.True(t, config.Record.SaveBodyToFile)
	assert.Equal(t, "nowhere", config.Record.SaveDir)
}
