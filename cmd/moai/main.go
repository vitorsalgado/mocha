package main

import (
	"bufio"
	"context"
	_ "embed"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/vitorsalgado/mocha/v3"
)

const (
	_dockerHostEnv    = "MOAI_DOCKER_HOST"
	_gitRepository    = "https://github.com/vitorsalgado/mocha"
	_usage            = "moai"
	_shortDescription = "Build Mock APIs in Go"
	_example          = `  moai --addr=:3000
  moai --proxy
  moai --proxy --record
  moai
`
)

var (
	//go:embed banner.txt
	_banner      string
	_description = fmt.Sprintf(`%s
Flexible HTTP mocking and expectations for Go. 
Supported mock file extensions: %s

For more information, visit: %s`,
		_banner,
		strings.Join(viper.SupportedExts, ", "),
		_gitRepository,
	)
)

func main() {
	rootCmd := &cobra.Command{
		Use:     _usage,
		Short:   _shortDescription,
		Long:    _description,
		Args:    cobra.MinimumNArgs(0),
		Example: _example,
		Run: func(cmd *cobra.Command, args []string) {
			conf := []mocha.Configurer{mocha.UseLocals()}

			_, exists := os.LookupEnv(_dockerHostEnv)
			if exists {
				conf = append(conf, &dockerConfigurer{})
			}

			m := mocha.New(conf...)

			ctx, cancel := signal.NotifyContext(m.Context(), os.Interrupt, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
			defer cancel()

			m.MustStart()

			fmt.Println(_banner)
			_ = m.PrintConfig(os.Stdin)

			readStdIn(ctx, inputs(bufio.NewReader(os.Stdin)))

			<-ctx.Done()

			m.Close()
		},
	}

	rootCmd.AddCommand(NewVersionCmd())

	err := rootCmd.Execute()
	if err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func inputs(reader *bufio.Reader) <-chan string {
	out := make(chan string)

	go func() {
		for {
			input, _ := reader.ReadString('\n')
			out <- input
		}
	}()

	return out
}

func readStdIn(ctx context.Context, in <-chan string) {
	go func() {
		for {
			select {
			case <-ctx.Done():
				return

			case input, ok := <-in:
				if !ok {
					return
				}

				println("done work " + input)
			}
		}
	}()
}
