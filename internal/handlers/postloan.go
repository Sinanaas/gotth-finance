package handlers

import (
	"github.com/Sinanaas/gotth-financial-tracker/internal/controllers"
	"github.com/gin-gonic/gin"
)

type PostLoanHandler struct {
	BC *controllers.BasicController
}

func NewPostLoanHandler(bc *controllers.BasicController) *PostLoanHandler {
	return &PostLoanHandler{BC: bc}
}

func (h *PostLoanHandler) ServeHTTP(c *gin.Context) {
	h.BC.CreateLoan(c)
}
