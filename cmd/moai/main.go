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

var (
	//go:embed banner.txt
	Banner     string
	Repository = "https://github.com/vitorsalgado/mocha"
	Usage      = "moai"
	Short      = "Build Mock APIs in Go"
	Example    = `  moai --addr=:3000
  moai --proxy
  moai --proxy --record
  moai
`
	Description = fmt.Sprintf(`%s
Flexible HTTP mocking and expectations for Go. 
Supported mock file extensions: %s

For more information, visit: %s`,
		Banner,
		strings.Join(viper.SupportedExts, ", "),
		Repository,
	)
)

func main() {
	rootCmd := &cobra.Command{
		Use:     Usage,
		Short:   Short,
		Long:    Description,
		Args:    cobra.MinimumNArgs(0),
		Example: Example,
		Run: func(cmd *cobra.Command, args []string) {
			m := mocha.New(mocha.UseLocalConfig(), mocha.UseFlags())

			ctx, cancel := signal.NotifyContext(m.Context(), os.Interrupt, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
			defer cancel()

			m.MustStart()

			fmt.Println(Banner)
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
