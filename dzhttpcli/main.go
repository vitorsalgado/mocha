package main

import (
	"context"
	_ "embed"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/vitorsalgado/mocha/v3/dzhttp"
)

const (
	_dockerHostEnv    = "DZ_DOCKER_HOST"
	_gitRepository    = "https://github.com/vitorsalgado/mocha"
	_usage            = "dz"
	_shortDescription = "Build Mock APIs in Go"
	_example          = `  dz --addr=:3000
  dz --proxy
  dz --proxy --record
  dz
  dz
`
)

var (
	//go:embed banner.txt
	_banner      string
	_description = fmt.Sprintf(`%s
Flexible HTTP mocking tool for Go.
Supported mock file extensions: %s

For more information, visit: %s`,
		_banner,
		strings.Join(viper.SupportedExts, ", "),
		_gitRepository,
	)

	locker sync.Mutex
)

func main() {
	ctx := context.Background()
	ctx, cancel := signal.NotifyContext(ctx, os.Interrupt, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	defer cancel()

	run(ctx) // should block
}

func run(ctx context.Context, custom ...dzhttp.Configurer) {
	configurers := make([]dzhttp.Configurer, 0)
	if len(custom) > 0 {
		configurers = append(configurers, custom...)
	}

	configurers = append(configurers, dzhttp.UseLocals())

	_, exists := os.LookupEnv(_dockerHostEnv)
	if exists {
		configurers = append(configurers, &dockerConfigurer{})
	}

	m := dzhttp.NewAPI(configurers...)

	rootCmd := &cobra.Command{
		Use:     _usage,
		Short:   _shortDescription,
		Long:    _description,
		Args:    cobra.MinimumNArgs(0),
		Example: _example,
		Run: func(cmd *cobra.Command, args []string) {
			m.MustStart()

			fmt.Println(_banner)
			fmt.Println(m.DescribeConfig())

			<-ctx.Done()

			m.Close()
		},
	}

	rootCmd.AddCommand(versionCmd())

	locker.Lock()
	defer locker.Unlock()

	err := rootCmd.Execute()
	if err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func versionCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "print version info",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println(dzhttp.Version)
		},
	}
}
