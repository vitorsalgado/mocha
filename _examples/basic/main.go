package main

import (
	"context"
	"fmt"
	_ "net/http/pprof"
	"os"
	"os/signal"
	"syscall"

	"github.com/vitorsalgado/mocha/v3"
	. "github.com/vitorsalgado/mocha/v3/matcher"
	"github.com/vitorsalgado/mocha/v3/misc"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	m := mocha.New(mocha.Setup().Addr(":8080"))
	m.MustStart()

	m.MustMock(mocha.
		Get(URLPath("/test")).
		Header(misc.HeaderAccept, Contain(misc.MIMETextHTML)).
		Header(misc.HeaderContentType, StrictEqual("test")).
		Header("any", All(Contain("test"), EqualIgnoreCase("dev"))).
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
