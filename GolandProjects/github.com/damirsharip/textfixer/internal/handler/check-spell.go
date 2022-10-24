package handler

import (
	"net/http"

	"textfixer/internal/models"

	"github.com/gin-gonic/gin"
)

func (h *Handler) checkTextSpelling(ctx *gin.Context) {

}

func (h *Handler) checkTextsSpelling(ctx *gin.Context) {
	var input models.TextsInput

	if err := ctx.ShouldBindJSON(&input); err != nil {
		//h.logger.Error(err)
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status": 400,
			"error":  "invalid input body provided",
		})
		return
	}

	res, err := h.services.CheckSpell(ctx, input)
	if err != nil {
		h.logger.Error(err)
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status": 500,
			"error":  "internal server error",
		})
		return
	}

	ctx.JSON(http.StatusOK, res)
}
