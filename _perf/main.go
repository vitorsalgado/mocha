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
	"github.com/vitorsalgado/mocha/v3/internal/header"
	"github.com/vitorsalgado/mocha/v3/internal/mimetype"
	. "github.com/vitorsalgado/mocha/v3/matcher"
	"github.com/vitorsalgado/mocha/v3/reply"
)

type Srv struct {
	h    http.Handler
	cfg  *mocha.Config
	info mocha.ServerInfo
}

func (s *Srv) Configure(config *mocha.Config, handler http.Handler) error {
	http.HandleFunc("/", handler.ServeHTTP)

	s.h = handler
	s.cfg = config

	return nil
}

func (s *Srv) Start() (mocha.ServerInfo, error) {
	go func() {
		if err := http.ListenAndServe(s.cfg.Addr, s.h); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalln(err)
		}
	}()

	return s.info, nil
}

func (s *Srv) StartTLS() (mocha.ServerInfo, error) {
	return s.info, nil
}

func (s *Srv) Close() error {
	return nil
}

func (s *Srv) Info() mocha.ServerInfo {
	return mocha.ServerInfo{URL: ""}
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

	m := mocha.New(mocha.NewConsoleNotifier(),
		mocha.Configure().
			HandlerDecorator(h).
			Server(&Srv{}).
			Addr(":8080").Build())
	m.Start()

	m.AddMocks(mocha.
		Get(URLPath("/test")).
		Header(header.Accept, Contain(mimetype.TextPlain)).
		Header("X-Scenario", Equal("1")).
		Reply(reply.OK().
			PlainText("ok").
			Header("X-Scenario-Result", "true")))

	m.AddMocks(mocha.
		Post(URLPath("/test")).
		Header(header.ContentType, Contain(mimetype.JSON)).
		Body(AllOf(
			JSONPath("active", Equal(true)),
			JSONPath("result", Equal("ok")))).
		Reply(reply.OK().
			JSON().
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
