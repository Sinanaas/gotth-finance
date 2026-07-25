package handlers

import (
	"github.com/Sinanaas/gotth-financial-tracker/internal/managers"
	"github.com/Sinanaas/gotth-financial-tracker/internal/templates"
	"github.com/Sinanaas/gotth-financial-tracker/internal/utils"
	"github.com/Sinanaas/gotth-financial-tracker/internal/validation"
	"github.com/gin-gonic/gin"
)

type DeleteLoanHandler struct {
	BM *managers.BasicManager
}

func NewDeleteLoanHandler(bm *managers.BasicManager) *DeleteLoanHandler {
	return &DeleteLoanHandler{BM: bm}
}

func (h *DeleteLoanHandler) ServeHTTP(ctx *gin.Context) {
	userId := utils.GetSessionUserID(ctx)

	loanId, err := validation.UUIDField(ctx.PostForm("LoanID"), "loan")
	if err != nil {
		swalError(ctx, err.Error())
		return
	}

	if err := h.BM.DeleteLoanById(loanId); err != nil {
		swalError(ctx, "Failed to delete loan")
		return
	}

	loans, err := h.BM.GetLoans(userId)
	if err != nil {
		swalError(ctx, "Failed to reload loans")
		return
	}
	ctx.Status(200)
	_ = templates.LoansPanel(loans).Render(ctx.Request.Context(), ctx.Writer)
}
