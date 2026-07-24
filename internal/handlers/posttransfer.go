package handlers

import (
	"net/http"

	"github.com/Sinanaas/gotth-financial-tracker/internal/controllers"
	"github.com/gin-gonic/gin"
)

type PostTransferHandler struct {
	BC *controllers.BasicController
}

func NewPostTransferHandler(bc *controllers.BasicController) *PostTransferHandler {
	return &PostTransferHandler{BC: bc}
}

func (h *PostTransferHandler) ServeHTTP(ctx *gin.Context) {
	if err := h.BC.CreateTransfer(ctx); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.Header("HX-Redirect", "/accounts")
}
