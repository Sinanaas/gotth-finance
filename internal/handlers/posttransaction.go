package handlers

import (
	"github.com/Sinanaas/gotth-financial-tracker/internal/controllers"
	"github.com/gin-gonic/gin"
)

type PostTransactionHandler struct {
	BC *controllers.BasicController
}

func NewPostTransactionHandler(bc *controllers.BasicController) *PostTransactionHandler {
	return &PostTransactionHandler{BC: bc}
}

func (h *PostTransactionHandler) ServeHTTP(c *gin.Context) {
	h.BC.CreateTransaction(c)
}
