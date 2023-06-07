package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/http/pprof"
	_ "net/http/pprof"
	"os"
	"os/signal"
	"path"
	"strings"
	"syscall"

	"github.com/vitorsalgado/mocha/v3"
	. "github.com/vitorsalgado/mocha/v3/matcher"
	mhttp2 "github.com/vitorsalgado/mocha/v3/mhttp"
	"github.com/vitorsalgado/mocha/v3/misc"
)

type Srv struct {
	h    http.Handler
	cfg  *mhttp2.Config
	info *mhttp2.ServerInfo
}

func (s *Srv) Setup(app *mhttp2.HTTPMockApp, handler http.Handler) error {
	http.HandleFunc("/", handler.ServeHTTP)

	s.h = handler
	s.cfg = app.Config()

	return nil
}

func (s *Srv) Start() error {
	go func() {
		if err := http.ListenAndServe(s.cfg.Addr, s.h); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalln(err)
		}
	}()

	return nil
}

func (s *Srv) StartTLS() error {
	return nil
}

func (s *Srv) Close() error {
	return nil
}

func (s *Srv) S() any {
	return nil
}

func (s *Srv) Info() *mhttp2.ServerInfo {
	return &mhttp2.ServerInfo{URL: ""}
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	wd, _ := os.Getwd()
	f, err := os.Open(path.Join(wd, "_perf", "res.json"))
	if err != nil {
		log.Fatalln(err)
	}

	defer f.Close()

	h := func(handler http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if strings.HasPrefix(r.URL.Path, "/debug/pprof/heap") {
				pprof.Index(w, r)
				return
			}

			if strings.HasPrefix(r.URL.Path, "/debug/pprof/profile") {
				pprof.Profile(w, r)
				return
			}

			if strings.HasPrefix(r.URL.Path, "/debug/pprof/cmdline") {
				pprof.Cmdline(w, r)
				return
			}

			if strings.HasPrefix(r.URL.Path, "/debug/pprof/symbol") {
				pprof.Symbol(w, r)
				return
			}

			if strings.HasPrefix(r.URL.Path, "/debug/pprof/trace") {
				pprof.Trace(w, r)
				return
			}

			if strings.HasPrefix(r.URL.Path, "/test") {
				handler.ServeHTTP(w, r)
				return
			}
		})
	}

	m := mocha.New(
		mhttp2.Setup().
			HandlerDecorator(h).
			Server(&Srv{}).
			Addr(":8080"))
	m.MustStart()

	m.MustMock(mhttp2.Get(URLPath("/test")).
		Header(misc.HeaderAccept, Contain(misc.MIMETextPlain)).
		Header("X-Scenario", StrictEqual("1")).
		Reply(mhttp2.OK().
			PlainText("ok").
			Header("X-Scenario-Result", "true")))

	m.MustMock(mhttp2.Post(URLPath("/test")).
		Header(misc.HeaderContentType, Contain(misc.MIMEApplicationJSON)).
		Body(All(
			JSONPath("active", StrictEqual(true)),
			JSONPath("result", StrictEqual("ok")))).
		Reply(mhttp2.OK().
			ContentType(misc.MIMEApplicationJSON).
			BodyReader(f)),
	)

	go func() {
		exit := make(chan os.Signal, 1)
		signal.Notify(exit, syscall.SIGINT, syscall.SIGTERM)

		<-exit

		cancel()
	}()

	fmt.Printf("runnning server for performance test: %s\n", m.URL())
	fmt.Printf("go to: %s\n", m.URL()+"/test")

	<-ctx.Done()

	m.Close()
}
