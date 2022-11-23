package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/vitorsalgado/mocha/v3"
	"github.com/vitorsalgado/mocha/v3/expect"
	"github.com/vitorsalgado/mocha/v3/internal/headerx"
	"github.com/vitorsalgado/mocha/v3/internal/mimetypex"
	"github.com/vitorsalgado/mocha/v3/reply"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	m := mocha.New(mocha.NewConsoleNotifier(), mocha.Configure().Context(ctx).Build())
	m.Start()

	m.AddMocks(mocha.
		Get(expect.URLPath("/test")).
		Header(headerx.Accept,
			expect.ToContain(mimetypex.TextHTML)).
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
