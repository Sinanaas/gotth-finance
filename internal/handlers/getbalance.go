package handlers

import (
	"github.com/Sinanaas/gotth-financial-tracker/internal/controllers"
	"github.com/gin-gonic/gin"
)

type GetBalanceHandler struct {
	BC *controllers.BasicController
}

func NewGetBalanceHandler(bc *controllers.BasicController) *GetBalanceHandler {
	return &GetBalanceHandler{BC: bc}
}

func (h *GetBalanceHandler) ServeHTTP(c *gin.Context) {
	h.BC.GetAccountBalance(c)
}
