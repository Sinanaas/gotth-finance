package handlers

import (
	"net/http"

	"github.com/Sinanaas/gotth-financial-tracker/internal/constants"
	"github.com/Sinanaas/gotth-financial-tracker/internal/controllers"
	"github.com/Sinanaas/gotth-financial-tracker/internal/templates"
	"github.com/Sinanaas/gotth-financial-tracker/internal/utils"
	"github.com/gin-gonic/gin"
)

type GetTransactionHandler struct {
	BC *controllers.BasicController
}

func NewGetTransaction(bc *controllers.BasicController) *GetTransactionHandler {
	return &GetTransactionHandler{BC: bc}
}

func (h *GetTransactionHandler) ServeHTTP(ctx *gin.Context) {
	cookie, _ := ctx.Cookie("access_token")
	userId := utils.GetSessionUserID(ctx)

	categories, err := h.BC.GetUserCategories(userId)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to load categories"})
		return
	}

	accounts, err := h.BC.GetUserAccounts(userId)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to load accounts"})
		return
	}

	transactions, total, err := h.BC.FilterTransactions(userId, "", "", "", "", "", -1, 1, 20)
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
