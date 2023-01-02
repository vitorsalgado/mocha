package main

import (
	"context"
	"fmt"
	_ "net/http/pprof"
	"os"
	"os/signal"
	"syscall"

	"github.com/vitorsalgado/mocha/v3"
	"github.com/vitorsalgado/mocha/v3/internal/header"
	"github.com/vitorsalgado/mocha/v3/internal/mimetype"
	. "github.com/vitorsalgado/mocha/v3/matcher"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	m := mocha.New(mocha.Configure().Addr(":8080"))
	m.MustStart()

	m.MustMock(mocha.
		Get(URLPath("/test")).
		Header(header.Accept, Contain(mimetype.TextHTML)).
		Header(header.ContentType, StrictEqual("test")).
		Header("any", AllOf(Contain("test"), EqualIgnoreCase("dev"))).
		Reply(mocha.OK().
			PlainText("hello world").
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

	m.Close()
}
