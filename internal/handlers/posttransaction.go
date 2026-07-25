package handlers

import (
	"github.com/Sinanaas/gotth-financial-tracker/internal/managers"
	"github.com/Sinanaas/gotth-financial-tracker/internal/models"
	"github.com/Sinanaas/gotth-financial-tracker/internal/templates"
	"github.com/Sinanaas/gotth-financial-tracker/internal/utils"
	"github.com/Sinanaas/gotth-financial-tracker/internal/validation"
	"github.com/gin-gonic/gin"
)

type PostTransactionHandler struct {
	BM *managers.BasicManager
}

func NewPostTransactionHandler(bm *managers.BasicManager) *PostTransactionHandler {
	return &PostTransactionHandler{BM: bm}
}

func (h *PostTransactionHandler) ServeHTTP(ctx *gin.Context) {
	userId := utils.GetSessionUserID(ctx)

	description, err := validation.Required(ctx.PostForm("Description"), "description")
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

	payload := models.TransactionRequest{
		Description: description,
		CategoryID:  categoryId,
		Amount:      amount,
		Date:        ctx.PostForm("Date"),
		Type:        txType,
		Account:     accountId,
		UserID:      userId,
	}
	if err := h.BM.CreateTransaction(payload); err != nil {
		swalError(ctx, err.Error())
		return
	}

	transactions, total, err := h.BM.FilterTransactions(userId, "", "", "", "", "", -1, 1, 20)
	if err != nil {
		swalError(ctx, "Failed to reload transactions")
		return
	}
	ctx.Status(200)
	_ = templates.TransactionListBody(transactions, total, 1, 20).Render(ctx.Request.Context(), ctx.Writer)
}
