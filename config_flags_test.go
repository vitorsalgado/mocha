package mocha

import (
	"testing"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func TestFlagsConfigurer_Apply(t *testing.T) {
	c := UseFlags().(*flagsConfigurer)
	assert.NoError(t, c.Apply(&Config{}))
}

func TestUseFlags(t *testing.T) {
	v := viper.New()
	c := UseFlags().(*flagsConfigurer)

	v.Set(FlagProxy, true)
	v.Set(FlagRecord, true)
	v.Set(FlagCORS, true)
	v.Set(FlagForwardTo, true)
	v.Set(FlagForwardTo, "https://www.example.org/")

	c.v = v

	assert.NoError(t, c.Apply(&Config{}))
}
