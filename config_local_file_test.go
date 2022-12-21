package mocha

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLocalConfig(t *testing.T) {
	c := UseLocalConfig().(*localConfigurer)

	assert.Equal(t, DefaultConfigFileName, c.filename)
	assert.Equal(t, DefaultConfigDirectories, c.paths)
}

func TestLocalConfig_ApplyExt(t *testing.T) {
	testCases := []struct {
		name     string
		filename string
	}{
		{".json", ".01_moairc.test.json"},
		{".yaml", ".02_moairc.test.yaml"},
		{".yml", ".02_moairc.test.yml"},
		{".properties", ".03_moairc.test.properties"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			cl := UseLocalConfigFrom(tc.filename, []string{"testdata"})
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
			assert.Equal(t, "https://proxy.org/test", config.Proxy.ProxyVia)
			assert.Equal(t, time.Duration(5000), config.Proxy.Timeout)

			assert.NotNil(t, config.Record)
			assert.Equal(t, []string{"header1", "header2"}, config.Record.RequestHeaders)
			assert.Equal(t, []string{"header3", "header4"}, config.Record.ResponseHeaders)
			assert.True(t, config.Record.Save)
			assert.True(t, config.Record.SaveBodyToFile)
			assert.Equal(t, "nowhere", config.Record.SaveDir)
		})
	}
}
