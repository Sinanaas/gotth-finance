package handlers

import (
	"github.com/Sinanaas/gotth-financial-tracker/internal/managers"
	"github.com/Sinanaas/gotth-financial-tracker/internal/templates"
	"github.com/Sinanaas/gotth-financial-tracker/internal/utils"
	"github.com/Sinanaas/gotth-financial-tracker/internal/validation"
	"github.com/gin-gonic/gin"
)

type PostFinishLoanHandler struct {
	BM *managers.BasicManager
}

func NewPostFinishLoanHandler(bm *managers.BasicManager) *PostFinishLoanHandler {
	return &PostFinishLoanHandler{BM: bm}
}

func (h *PostFinishLoanHandler) ServeHTTP(ctx *gin.Context) {
	userId := utils.GetSessionUserID(ctx)

	loanId, err := validation.UUIDField(ctx.PostForm("LoanID"), "loan")
	if err != nil {
		swalError(ctx, err.Error())
		return
	}

	if err := h.BM.FinishLoan(loanId, userId); err != nil {
		swalError(ctx, err.Error())
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
