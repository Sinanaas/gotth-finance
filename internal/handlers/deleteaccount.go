package handlers

import (
	"github.com/Sinanaas/gotth-financial-tracker/internal/controllers"
	"github.com/gin-gonic/gin"
)

type DeleteAccountHandler struct {
	BC *controllers.BasicController
}

func NewDeleteAccountHandler(bc *controllers.BasicController) *DeleteAccountHandler {
	return &DeleteAccountHandler{BC: bc}
}

func (h *DeleteAccountHandler) ServeHTTP(ctx *gin.Context) {
	h.BC.DeleteAccountById(ctx)
}
