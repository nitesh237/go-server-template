package http

import (
	"context"
	"fmt"
	"net"
	"net/http"

	"go.uber.org/fx"
	"go.uber.org/zap"

	"github.com/nitesh237/go-server-template/pkg/cfg"
	"github.com/nitesh237/go-server-template/pkg/log"
)

var (
	FxHttpServerModule = fx.Module("http-server",
		fx.Provide(
			NewHTTPServer,
		),
		fx.Invoke(func(s *http.Server) {}),
	)
)

func NewHTTPServer(lc fx.Lifecycle, appCong *cfg.Application, handler http.Handler, logger log.Logger) *http.Server {
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", appCong.ServerPorts.HttpPort),
		Handler: handler,
	}
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			ln, err := net.Listen("tcp", srv.Addr)
			if err != nil {
				return err
			}
			logger.InfoNoCtx("Starting HTTP server", zap.String("addr", srv.Addr))
			go func() {
				serverErr := srv.Serve(ln)
				if serverErr != nil && serverErr != http.ErrServerClosed {
					logger.ErrorNoCtx("HTTP server error", zap.Error(serverErr))
				}
			}()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			return srv.Shutdown(ctx)
		},
	})
	return srv
}
