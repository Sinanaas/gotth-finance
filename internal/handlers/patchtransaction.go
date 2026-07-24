package handlers

import (
	"net/http"

	"github.com/Sinanaas/gotth-financial-tracker/internal/controllers"
	"github.com/gin-gonic/gin"
)

type PatchTransactionHandler struct {
	BC *controllers.BasicController
}

func NewPatchTransactionHandler(bc *controllers.BasicController) *PatchTransactionHandler {
	return &PatchTransactionHandler{BC: bc}
}

func (h *PatchTransactionHandler) ServeHTTP(ctx *gin.Context) {
	if err := h.BC.UpdateTransaction(ctx); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.Header("HX-Redirect", "/transaction")
}
