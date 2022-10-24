package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *Handler) healthCheck(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{
		"status": 200,
	})
}
