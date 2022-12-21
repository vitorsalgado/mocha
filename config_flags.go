package mocha

import (
	"fmt"
	"net/url"
	"os"
	"time"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

// Configuration flags.
const (
	FlagConfigFilename        = "config"
	FlagConfigFilenameShort   = "c"
	FlagName                  = "name"
	FlagNameShort             = "n"
	FlagUseHTTPS              = "https"
	FlagUseHTTPSShort         = "s"
	FlagAddr                  = "addr"
	FlagAddrShort             = "a"
	FlagVerbose               = "verbose"
	FlagVerboseShort          = "v"
	FlagInfo                  = "info"
	FlagSilent                = "silent"
	FlagGlob                  = "glob"
	FlagGlobShort             = "g"
	FlagRecord                = "record"
	FlagRecordDir             = "record-dir"
	FlagRecordSaveBodyToFile  = "record-save-body-file"
	FlagRecordRequestHeaders  = "record-request-headers"
	FlagRecordResponseHeaders = "record-response-headers"
	FlagProxy                 = "proxy"
	FlagProxyShort            = "p"
	FlagProxyVia              = "proxy-via"
	FlagProxyTimeout          = "proxy-timeout"
	FlagCORS                  = "cors"
	FlagForwardTo             = "forward-to"
)

var _ Configurer = (*flagsConfigurer)(nil)

type flagsConfigurer struct {
	v *viper.Viper
}

func (c *flagsConfigurer) Apply(conf *Config) error {
	var (
		name                  string
		configFilename        string
		useHTTPS              bool
		addr                  string
		verbose               bool
		info                  bool
		silent                bool
		directories           []string
		rec                   bool
		recordDir             string
		recordSaveBodyToFiles bool
		recordRequestHeaders  []string
		recordResponseHeaders []string
		proxy                 bool
		proxyVia              string
		proxyTimeout          int64
		cors                  bool
		forwardTo             string
	)

	var (
		defRec = defaultRecordConfig()
	)

	c.v.SetDefault(FlagName, conf.Name)
	c.v.SetDefault(FlagVerbose, conf.LogLevel == LogVerbose)
	c.v.SetDefault(FlagAddr, conf.Addr)
	c.v.SetDefault(FlagGlob, conf.Directories)

	pp := pflag.NewFlagSet("", pflag.ContinueOnError)

	pp.StringVarP(&configFilename, FlagConfigFilename, FlagConfigFilenameShort, "", "Use a custom configuration file. Commandline arguments take precedence over file values.")
	pp.StringVarP(&name, FlagName, FlagNameShort, "", "Mock server name.")
	pp.BoolVarP(&useHTTPS, FlagUseHTTPS, FlagUseHTTPSShort, false, "Enable HTTPS.")
	pp.StringVarP(&addr, FlagAddr, FlagAddrShort, "", "Server address. Usage: <host>:<port>, :<port>, <port>. If no value is set, it will an auto generated one.")
	pp.BoolVarP(&verbose, FlagVerbose, FlagVerboseShort, false, "Verbose logs")
	pp.BoolVar(&info, FlagInfo, false, "Informative logs")
	pp.BoolVar(&silent, FlagSilent, false, "Minimum logs")
	pp.StringSliceVarP(&directories, FlagGlob, FlagGlobShort, []string{}, "Mock search glob patterns. Example: testdata/*mock.json,testdata/*mock.yaml")

	pp.BoolVar(&rec, FlagRecord, false, "Enable mock recording.")
	pp.StringVar(&recordDir, FlagRecordDir, "testdata/", "Recorded mocks directory.")
	pp.BoolVar(&recordSaveBodyToFiles, FlagRecordSaveBodyToFile, false, "Save recorded body to separate file, using \"body_file\" field.")
	pp.StringSliceVar(&recordRequestHeaders, FlagRecordRequestHeaders, defRec.RequestHeaders, "Request headers to be recorded.")
	pp.StringSliceVar(&recordResponseHeaders, FlagRecordResponseHeaders, defRec.ResponseHeaders, "Response headers to be recorded.")

	pp.BoolVarP(&proxy, FlagProxy, FlagProxyShort, false, "Enable reverse proxy.")
	pp.StringVar(&proxyVia, FlagProxyVia, "", "Route proxy requests to another proxy.")
	pp.Int64Var(&proxyTimeout, FlagProxyTimeout, _defaultProxyConfig.Timeout.Milliseconds(), "Proxy timeout in milliseconds.")

	pp.BoolVar(&cors, FlagCORS, false, "Enable CORS.")

	pp.StringVar(&forwardTo, FlagForwardTo, "", "Forward requests to the given URL.")

	err := pp.Parse(os.Args)
	if err != nil {
		return err
	}

	err = c.v.BindPFlags(pp)
	if err != nil {
		return err
	}

	conf.UseHTTPS = c.v.GetBool(FlagUseHTTPS)
	conf.Addr = c.v.GetString(FlagAddr)
	conf.Name = c.v.GetString(FlagName)
	conf.Directories = c.v.GetStringSlice(FlagGlob)

	if c.v.GetBool(FlagProxy) {
		conf.Proxy = &ProxyConfig{
			ProxyVia: c.v.GetString(FlagProxyVia),
			Timeout:  time.Duration(c.v.GetInt64(FlagProxyTimeout)),
		}
	}

	if c.v.GetBool(FlagRecord) {
		conf.Record = defRec
		conf.Record.SaveDir = c.v.GetString(FlagRecordDir)
		conf.Record.SaveBodyToFile = c.v.GetBool(FlagRecordSaveBodyToFile)
		conf.Record.RequestHeaders = c.v.GetStringSlice(FlagRecordRequestHeaders)
		conf.Record.ResponseHeaders = c.v.GetStringSlice(FlagRecordResponseHeaders)
	}

	if c.v.GetBool(FlagCORS) {
		conf.CORS = &_defaultCORSConfig
	}

	if c.v.IsSet(FlagForwardTo) {
		target := c.v.GetString(FlagForwardTo)
		if target == "" {
			return fmt.Errorf(`flag "--%s" is required`, FlagForwardTo)
		}

		targetURL, err := url.Parse(target)
		if err != nil {
			return fmt.Errorf(`flag "--%s" must contain a valid URL. %w`, FlagForwardTo, err)
		}

		conf.Forward = &ForwardConfig{
			Target: targetURL.String(),
		}
	}

	return nil
}

// UseFlags configures the mock server using command-line flags.
func UseFlags() Configurer {
	return &flagsConfigurer{v: viper.New()}
}
