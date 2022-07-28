package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/vitorsalgado/mocha/v2"
	"github.com/vitorsalgado/mocha/v2/expect"
	"github.com/vitorsalgado/mocha/v2/internal/headers"
	"github.com/vitorsalgado/mocha/v2/internal/mimetypes"
	"github.com/vitorsalgado/mocha/v2/reply"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	m := mocha.New(mocha.NewConsoleNotifier(), mocha.Configure().Context(ctx).Build())
	m.Start()

	m.AddMocks(mocha.
		Get(expect.URLPath("/test")).
		Header(headers.Accept,
			expect.ToContain(mimetypes.TextHTML)).
		Reply(reply.OK().
			BodyString("hello world").
			Header("x-basic", "true")))

	go func() {
		exit := make(chan os.Signal, 1)
		signal.Notify(exit, syscall.SIGINT, syscall.SIGTERM)

		<-exit

		cancel()
	}()

	fmt.Printf("runnning basic example on: %s\n", m.URL())
	fmt.Printf("go to: %s\n", m.URL()+"/test")

	<-ctx.Done()
}
