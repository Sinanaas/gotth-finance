package handlers

import (
	"github.com/Sinanaas/gotth-financial-tracker/internal/managers"
	"github.com/gin-gonic/gin"
)

type DeleteTransactionHandler struct {
	BM *managers.BasicManager
}

func NewDeleteTransactionHandler(bm *managers.BasicManager) *DeleteTransactionHandler {
	return &DeleteTransactionHandler{BM: bm}
}

func (h *DeleteTransactionHandler) ServeHTTP(ctx *gin.Context) {
	transactionId := ctx.PostForm("TransactionID")
	if transactionId == "" {
		ctx.JSON(400, gin.H{"error": "Transaction ID is missing"})
		return
	}

	err := h.BM.DeleteTransactionById(transactionId)
	if err != nil {
		ctx.JSON(500, gin.H{"error": "Failed to delete transaction"})
		return
	}

	accountId := ctx.PostForm("AccountID")

	err = h.BM.RecalculateAccountBalance(accountId)
	if err != nil {
		ctx.JSON(500, gin.H{"error": "Failed to recalculate balance"})
		return
	}

	ctx.Header("HX-Redirect", "/transaction")
}
