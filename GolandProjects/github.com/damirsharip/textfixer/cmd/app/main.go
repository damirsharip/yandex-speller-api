package main

import (
	"context"

	"go.uber.org/zap"
	handler "textfixer/internal/handler"
	"textfixer/internal/server"
	"textfixer/pkg/speller"

	"github.com/gin-gonic/gin"
)

func main() {
	ctx := context.Background()
	gin.SetMode(gin.ReleaseMode)

	logger, _ := zap.NewProduction()
	defer logger.Sync() // flushes buffer, if any
	sugar := logger.Sugar()

	// initializing docusign service receiving access token for future use
	err := speller.NewClient()
	if err != nil {
		//logger.FatalKV(initCtx, "failed while docusign init: %s", err.Error())
	}

	//srvc := service.NewService()
	h := handler.NewHandler(sugar)

	// setting up routes
	router := h.InitRoutes()

	// initializing server
	srv := server.NewServer(sugar, router)

	if err := srv.Run(ctx); err != nil {
		sugar.Fatal("errored while starting the server", err.Error())
	}
}
