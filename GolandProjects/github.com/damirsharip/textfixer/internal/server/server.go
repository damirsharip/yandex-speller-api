package server

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"sync/atomic"
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

	isReady := &atomic.Value{}
	isReady.Store(false)

	listener, err = net.Listen("tcp", fmt.Sprintf("%s:%d", "localhost", 8000))

	if err != nil {
		//srv.logger.Fatal(err)
		//srv.logger.FatalKV(ctx, err.Error())
	}

	defer listener.Close()

	srv.httpServer = &http.Server{
		Handler: srv.handler,
		//cors(srv.handler, cfg.Rest.AllowedCORSOrigins),
		MaxHeaderBytes: 40 << 20, // 1 MB
		WriteTimeout:   15 * time.Second,
		ReadTimeout:    15 * time.Second,
	}
	srv.logger.Infof("starting listening on: %s:%d", "localhost", 8000)
	go func() {
		if err = srv.httpServer.Serve(listener); err != nil && !errors.Is(err, http.ErrServerClosed) {
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
		//logger.ErrorKV(ctx, fmt.Sprintf("Server.Shutdown(): %s", err.Error()))
		return err
	}

	return nil
}

//
//func cors(h http.Handler, allowedOrigins []string) http.Handler {
//	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
//		providedOrigin := r.Header.Get("Origin")
//		matches := false
//		for _, allowedOrigin := range allowedOrigins {
//			if providedOrigin == allowedOrigin {
//				matches = true
//				break
//			}
//		}
//
//		if matches {
//			w.Header().Set("Access-Control-Allow-Origin", providedOrigin)
//			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PATCH, DELETE")
//			w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, Authorization, ResponseType, X-Request-ID, x-payload-digest, x-authorization-digest")
//		}
//		if r.Method == "OPTIONS" {
//			return
//		}
//		h.ServeHTTP(w, r)
//	})
//}
