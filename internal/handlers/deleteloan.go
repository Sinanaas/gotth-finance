package handlers

import (
	"github.com/Sinanaas/gotth-financial-tracker/internal/controllers"
	"github.com/gin-gonic/gin"
)

type DeleteLoanHandler struct {
	BC *controllers.BasicController
}

func NewDeleteLoanHandler(bc *controllers.BasicController) *DeleteLoanHandler {
	return &DeleteLoanHandler{BC: bc}
}

func (h *DeleteLoanHandler) ServeHTTP(ctx *gin.Context) {
	h.BC.DeleteLoanById(ctx)
}
