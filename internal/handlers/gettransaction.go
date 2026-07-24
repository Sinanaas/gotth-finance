package handlers

import (
	"net/http"

	"github.com/Sinanaas/gotth-financial-tracker/internal/constants"
	"github.com/Sinanaas/gotth-financial-tracker/internal/managers"
	"github.com/Sinanaas/gotth-financial-tracker/internal/templates"
	"github.com/Sinanaas/gotth-financial-tracker/internal/utils"
	"github.com/gin-gonic/gin"
)

type GetTransactionHandler struct {
	BM *managers.BasicManager
}

func NewGetTransaction(bm *managers.BasicManager) *GetTransactionHandler {
	return &GetTransactionHandler{BM: bm}
}

func (h *GetTransactionHandler) ServeHTTP(ctx *gin.Context) {
	cookie, _ := ctx.Cookie("access_token")
	userId := utils.GetSessionUserID(ctx)

	categories, err := h.BM.GetUserCategories(userId)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to load categories"})
		return
	}

	accounts, err := h.BM.GetUserAccounts(userId)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to load accounts"})
		return
	}

	transactions, total, err := h.BM.FilterTransactions(userId, "", "", "", "", "", -1, 1, 20)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to load transactions"})
		return
	}

	var txType constants.TransactionType
	transactionTypeArray := txType.ToArrayString()

	c := templates.Transaction(categories, transactions, accounts, transactionTypeArray, total)
	err = templates.Layout(c, cookie).Render(ctx.Request.Context(), ctx.Writer)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
}
