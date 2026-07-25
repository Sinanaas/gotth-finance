package handlers

import (
	"github.com/Sinanaas/gotth-financial-tracker/internal/managers"
	"github.com/Sinanaas/gotth-financial-tracker/internal/templates"
	"github.com/Sinanaas/gotth-financial-tracker/internal/utils"
	"github.com/Sinanaas/gotth-financial-tracker/internal/validation"
	"github.com/gin-gonic/gin"
)

type DeleteTransactionHandler struct {
	BM *managers.BasicManager
}

func NewDeleteTransactionHandler(bm *managers.BasicManager) *DeleteTransactionHandler {
	return &DeleteTransactionHandler{BM: bm}
}

func (h *DeleteTransactionHandler) ServeHTTP(ctx *gin.Context) {
	userId := utils.GetSessionUserID(ctx)

	transactionId, err := validation.UUIDField(ctx.PostForm("TransactionID"), "transaction")
	if err != nil {
		swalError(ctx, err.Error())
		return
	}
	accountId, err := validation.UUIDField(ctx.PostForm("AccountID"), "account")
	if err != nil {
		swalError(ctx, err.Error())
		return
	}

	if err := h.BM.DeleteTransactionById(transactionId); err != nil {
		swalError(ctx, "Failed to delete transaction")
		return
	}
	if err := h.BM.RecalculateAccountBalance(accountId); err != nil {
		swalError(ctx, "Failed to recalculate balance")
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
