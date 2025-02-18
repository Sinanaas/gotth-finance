package handlers

import (
	"github.com/Sinanaas/gotth-financial-tracker/internal/controllers"
	"github.com/gin-gonic/gin"
)

type DeleteRecurringHandler struct {
	BC *controllers.BasicController
}

func NewDeleteRecurringHandler(bc *controllers.BasicController) *DeleteRecurringHandler {
	return &DeleteRecurringHandler{BC: bc}
}

func (h *DeleteRecurringHandler) ServeHTTP(ctx *gin.Context) {
	h.BC.DeleteRecurringById(ctx)
}
