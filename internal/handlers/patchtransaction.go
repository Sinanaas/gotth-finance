package handlers

import (
	"github.com/Sinanaas/gotth-financial-tracker/internal/managers"
	"github.com/Sinanaas/gotth-financial-tracker/internal/models"
	"github.com/Sinanaas/gotth-financial-tracker/internal/templates"
	"github.com/Sinanaas/gotth-financial-tracker/internal/utils"
	"github.com/Sinanaas/gotth-financial-tracker/internal/validation"
	"github.com/gin-gonic/gin"
)

type PatchTransactionHandler struct {
	BM *managers.BasicManager
}

func NewPatchTransactionHandler(bm *managers.BasicManager) *PatchTransactionHandler {
	return &PatchTransactionHandler{BM: bm}
}

func (h *PatchTransactionHandler) ServeHTTP(ctx *gin.Context) {
	userId := utils.GetSessionUserID(ctx)

	transactionId, err := validation.UUIDField(ctx.Param("id"), "transaction")
	if err != nil {
		swalError(ctx, err.Error())
		return
	}
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
	newAccountId, err := validation.UUIDField(ctx.PostForm("Account"), "account")
	if err != nil {
		swalError(ctx, err.Error())
		return
	}
	oldAccountId, err := validation.UUIDField(ctx.PostForm("OldAccountID"), "account")
	if err != nil {
		swalError(ctx, err.Error())
		return
	}
	if _, err := validation.Date(ctx.PostForm("Date")); err != nil {
		swalError(ctx, err.Error())
		return
	}

	payload := models.TransactionUpdateRequest{
		ID:           transactionId,
		Description:  description,
		CategoryID:   categoryId,
		OldAccountID: oldAccountId,
		NewAccountID: newAccountId,
		Amount:       amount,
		Date:         ctx.PostForm("Date"),
		Type:         txType,
		UserID:       userId,
	}
	if err := h.BM.UpdateTransaction(payload); err != nil {
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
