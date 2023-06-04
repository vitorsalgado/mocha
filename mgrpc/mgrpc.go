package mgrpc

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"reflect"

	"github.com/rs/zerolog"
	"google.golang.org/grpc"

	"github.com/vitorsalgado/mocha/v3/foundation"
)

type GRPCMockApp struct {
	*foundation.BaseApp[*GRPCMock, *GRPCMockApp]

	ctx     context.Context
	cancel  context.CancelFunc
	server  *grpc.Server
	config  *Config
	logger  *zerolog.Logger
	addr    string
	storage *foundation.MockStore[*GRPCMock]
}

type ServerInfo struct {
	Addr string
}

func NewGRPC(config ...foundation.Configurer[*Config]) *GRPCMockApp {
	app := &GRPCMockApp{}
	conf := defaultConfig()

	for i, configurer := range config {
		err := configurer.Apply(conf)
		if err != nil {
			panic(fmt.Errorf(
				"server: error applying configuration at index %d with type %v\n%w",
				i,
				reflect.TypeOf(configurer),
				err,
			))
		}
	}

	store := foundation.NewStore[*GRPCMock]()
	ctx, cancel := context.WithCancel(conf.Context)
	in := &Interceptors{app: app}
	srv := grpc.NewServer(
		grpc.UnaryInterceptor(in.UnaryInterceptor),
		grpc.StreamInterceptor(in.StreamInterceptor),
	)

	if conf.ServiceDesc != nil && conf.Service != nil {
		srv.RegisterService(conf.ServiceDesc, conf.Service)
	}

	app.BaseApp = foundation.NewBaseApp(app, store)
	app.storage = store
	app.ctx = ctx
	app.cancel = cancel
	app.config = conf
	app.server = srv

	return app
}

func NewGRPCWithT(t foundation.TestingT, config ...foundation.Configurer[*Config]) *GRPCMockApp {
	app := NewGRPC(config...)
	t.Cleanup(app.Close)

	return app
}

// Start starts the mock server.
func (app *GRPCMockApp) Start() error {
	lis, err := net.Listen("tcp", app.config.Addr)
	if err != nil {
		return fmt.Errorf("server: Start: failed to create listener: %w", err)
	}

	app.addr = lis.Addr().String()

	go func() {
		if err = app.server.Serve(lis); err != nil && err != http.ErrServerClosed {
			app.logger.Warn().Err(err).Msg("server: Start: error listening")
		}
	}()

	go func() {
		<-app.ctx.Done()
		app.server.GracefulStop()
	}()

	return nil
}

// MustStart starts the mock server.
// It fails immediately if any error occurs.
func (app *GRPCMockApp) MustStart() {
	err := app.Start()
	if err != nil {
		panic(err)
	}
}

func (app *GRPCMockApp) Server() *grpc.Server {
	return app.server
}

func (app *GRPCMockApp) Addr() string {
	return app.addr
}

func (app *GRPCMockApp) Close() {
	app.cancel()
}