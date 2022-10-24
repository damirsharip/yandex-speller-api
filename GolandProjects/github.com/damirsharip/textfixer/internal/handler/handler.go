package handler

import (
	"go.uber.org/zap"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	logger *zap.SugaredLogger
	//ctx      context.Context
	//services service.Service
	srv *gin.Engine
}

//services *service.Service
func NewHandler(logger *zap.SugaredLogger) *Handler {
	return &Handler{
		//ctx:      ctx,
		//services: *services,
		logger: logger,
	}
}

func (h *Handler) InitRoutes() *gin.Engine {
	router := gin.New()

	router.Use(gin.Recovery())

	router.GET("/", h.healthCheck)

	api := router.Group("/api")
	{
		v := api.Group("/v0")
		{
			v.GET("/healthz", h.healthCheck)
			v.POST("/check-texts-spelling", h.checkTextsSpelling)
			v.POST("/check-text-spelling", h.checkTextSpelling)
		}
	}

	return router
}
