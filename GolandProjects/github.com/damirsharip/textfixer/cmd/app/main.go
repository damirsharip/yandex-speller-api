package main

import (
	"context"

	"go.uber.org/zap"
	handler "textfixer/internal/handler"
	"textfixer/internal/server"
	"textfixer/internal/service"
)

func main() {
	ctx := context.Background()

	logger, _ := zap.NewProduction()
	defer logger.Sync() // flushes buffer, if any
	sugar := logger.Sugar()

	srvc := service.NewService()
	h := handler.NewHandler(sugar, srvc)

	// setting up routes
	router := h.InitRoutes()

	// initializing server
	srv := server.NewServer(sugar, router)

	if err := srv.Run(ctx); err != nil {
		sugar.Fatal("errored while starting the server", err.Error())
	}
}
