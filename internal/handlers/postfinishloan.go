package handlers

import (
	"github.com/Sinanaas/gotth-financial-tracker/internal/controllers"
	"github.com/gin-gonic/gin"
)

type PostFinishLoanHandler struct {
	BC *controllers.BasicController
}

func NewPostFinishLoanHandler(bc *controllers.BasicController) *PostFinishLoanHandler {
	return &PostFinishLoanHandler{BC: bc}
}

func (h *PostFinishLoanHandler) ServeHTTP(c *gin.Context) {
	h.BC.FinishLoan(c)
}
