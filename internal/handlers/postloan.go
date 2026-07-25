package handlers

import (
	"github.com/Sinanaas/gotth-financial-tracker/internal/managers"
	"github.com/Sinanaas/gotth-financial-tracker/internal/models"
	"github.com/Sinanaas/gotth-financial-tracker/internal/templates"
	"github.com/Sinanaas/gotth-financial-tracker/internal/utils"
	"github.com/Sinanaas/gotth-financial-tracker/internal/validation"
	"github.com/gin-gonic/gin"
)

type PostLoanHandler struct {
	BM *managers.BasicManager
}

func NewPostLoanHandler(bm *managers.BasicManager) *PostLoanHandler {
	return &PostLoanHandler{BM: bm}
}

func (h *PostLoanHandler) ServeHTTP(ctx *gin.Context) {
	userId := utils.GetSessionUserID(ctx)

	toWhom, err := validation.Required(ctx.PostForm("Towhom"), "person or party")
	if err != nil {
		swalError(ctx, err.Error())
		return
	}
	amount, err := validation.Amount(ctx.PostForm("Amount"))
	if err != nil {
		swalError(ctx, err.Error())
		return
	}
	txType, err := validation.IntField(ctx.PostForm("Type"), "type")
	if err != nil {
		swalError(ctx, err.Error())
		return
	}
	categoryId, err := validation.UUIDField(ctx.PostForm("Category"), "category")
	if err != nil {
		swalError(ctx, err.Error())
		return
	}
	accountId, err := validation.UUIDField(ctx.PostForm("Account"), "account")
	if err != nil {
		swalError(ctx, err.Error())
		return
	}
	if _, err := validation.Date(ctx.PostForm("Date")); err != nil {
		swalError(ctx, err.Error())
		return
	}

	payload := models.LoanRequest{
		Description:     ctx.PostForm("Description"),
		CategoryID:      categoryId,
		ToWhom:          toWhom,
		Status:          false,
		LoanDate:        ctx.PostForm("Date"),
		Amount:          amount,
		TransactionType: txType,
		AccountID:       accountId,
		UserID:          userId,
	}
	if err := h.BM.CreateLoan(payload); err != nil {
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
