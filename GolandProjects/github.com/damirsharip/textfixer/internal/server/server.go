package server

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"go.uber.org/zap"
)

type Server struct {
	handler    http.Handler
	logger     *zap.SugaredLogger
	httpServer *http.Server
}

func NewServer(logger *zap.SugaredLogger, h http.Handler) *Server {
	return &Server{
		handler: h,
		logger:  logger,
	}
}

func (srv *Server) Run(ctx context.Context) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	var listener net.Listener
	var err error

	listener, err = net.Listen("tcp", fmt.Sprintf("%s:%d", "localhost", 8000))

	if err != nil {
		srv.logger.Fatal(err)
	}

	defer func() {
		err = listener.Close()
	}()

	srv.httpServer = &http.Server{
		Handler:        srv.handler,
		MaxHeaderBytes: 1 << 20, // 1 MB
		WriteTimeout:   15 * time.Second,
		ReadTimeout:    15 * time.Second,
	}
	srv.logger.Infof("starting listening on: %s:%d", "localhost", 8000)
	go func() {
		if err = srv.httpServer.Serve(listener); err != nil && !errors.Is(err, http.ErrServerClosed) {
			srv.logger.Error(err)
			cancel()
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	select {
	case v := <-quit:
		srv.logger.Info(ctx, fmt.Sprintf("signal.Notify: %v", v))
	case done := <-ctx.Done():
		srv.logger.Info(ctx, fmt.Sprintf("ctx.Done: %v", done))
	}
	ctx, finalCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer finalCancel()

	if err = srv.httpServer.Shutdown(ctx); err != nil {
		srv.logger.Error(fmt.Errorf("Server.Shutdown(): %s", err.Error()))
		return err
	}

	return nil
}
