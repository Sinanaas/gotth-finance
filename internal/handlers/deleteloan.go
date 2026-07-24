package handlers

import (
	"github.com/Sinanaas/gotth-financial-tracker/internal/managers"
	"github.com/gin-gonic/gin"
)

type DeleteLoanHandler struct {
	BM *managers.BasicManager
}

func NewDeleteLoanHandler(bm *managers.BasicManager) *DeleteLoanHandler {
	return &DeleteLoanHandler{BM: bm}
}

func (h *DeleteLoanHandler) ServeHTTP(ctx *gin.Context) {
	loanId := ctx.PostForm("LoanID")
	if loanId == "" {
		ctx.JSON(400, gin.H{"error": "Loan ID is missing"})
		return
	}

	err := h.BM.DeleteLoanById(loanId)
	if err != nil {
		ctx.JSON(500, gin.H{"error": "Failed to delete loan"})
		return
	}
	ctx.Header("HX-Redirect", "/loans")
}
