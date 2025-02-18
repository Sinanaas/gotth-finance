package handlers

import (
	"github.com/Sinanaas/gotth-financial-tracker/internal/controllers"
	"github.com/gin-gonic/gin"
)

type DeleteTransactionHandler struct {
	BC *controllers.BasicController
}

func NewDeleteTransactionHandler(bc *controllers.BasicController) *DeleteTransactionHandler {
	return &DeleteTransactionHandler{BC: bc}
}

func (h *DeleteTransactionHandler) ServeHTTP(ctx *gin.Context) {
	h.BC.DeleteTransactionById(ctx)
}
