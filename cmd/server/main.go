package main

import (
	"audita-api-gateway/options"
	"context"
	"net"
	"net/http"
	httpCaller "net/http"
	"os"
	"sync/atomic"
	"syscall"
	"time"

	"audita-api-gateway/pkg/endpoint"
	"audita-api-gateway/pkg/service"
	"audita-api-gateway/transports"

	"github.com/audita-bids/private-kit/middlewares"
	"github.com/audita-bids/private-kit/pkg/lib"
	"github.com/go-kit/kit/log/level"
	"github.com/oklog/run"
)

func main() {
	instanceCfg := options.NewCfg()
	cfg := options.HandleCfg(instanceCfg)

	logger := lib.SetupLogger(cfg.Debug)

	var (
		svc         = service.NewService(logger)
		endpoints   = endpoint.NewEndpointSetup(svc, logger)
		httpHandler = transports.NewHTTPServer(*endpoints, logger)

		httpServer *httpCaller.Server
	)

	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	f := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Company-Id")

			if r.Method == http.MethodOptions {
				w.WriteHeader(http.StatusOK)
				return
			}

			next.ServeHTTP(w, r)
		})
	}

	ready := new(atomic.Bool)

	var g run.Group
	{
		httpListener, err := net.Listen("tcp", cfg.HttpAddr)
		if err != nil {
			level.Error(logger).Log("msg", "failed to listen on http address", "err", err)
			os.Exit(1)
		}

		strip := http.StripPrefix("/api", httpHandler)

		httpServer = &httpCaller.Server{
			Handler:           f(strip),
			ReadHeaderTimeout: 5 * time.Second,
			ReadTimeout:       30 * time.Second,
			WriteTimeout:      120 * time.Second,
			IdleTimeout:       120 * time.Second,
		}

		g.Add(func() error {
			return httpServer.Serve(httpListener)
		}, func(error) {
			ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
			defer cancel()

			if err := httpServer.Shutdown(ctx); err != nil {
				level.Error(logger).Log("msg", "failed to shutdown http server", "err", err)
			}
		})
	}
	{
		promListener, err := net.Listen("tcp", cfg.PromAddr)
		if err != nil {
			level.Error(logger).Log("msg", "failed to listen on prometheus address", "err", err)
			os.Exit(1)
		}

		config := middlewares.MetricsConfig{
			Logger:         logger,
			EnableEndpoint: true,
			EnableHTTP:     true,
			ServiceName:    "contracts",
			Ready:          ready,
		}

		srv := middlewares.NewMetricsServer(config, cfg.PromAddr)

		g.Add(func() error {
			level.Info(logger).Log(
				"msg", "prometheus server started",
				"addr", cfg.PromAddr,
			)

			return srv.Serve(promListener)
		}, func(error) {
			ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
			defer cancel()

			if err := srv.Shutdown(ctx); err != nil {
				level.Error(logger).Log("msg", "failed to shutdown prometheus server", "err", err)
			}
		})
	}

	{
		g.Add(run.SignalHandler(ctx, syscall.SIGINT, syscall.SIGTERM))
	}

	ready.Store(true)

	level.Info(logger).Log("msg", "starting servers")
	if err := g.Run(); err != nil {
		level.Error(logger).Log("msg", "servers failed", "err", err)
		os.Exit(1)
	}
}
