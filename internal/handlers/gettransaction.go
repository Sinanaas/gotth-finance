package handlers

import (
	"github.com/Sinanaas/gotth-financial-tracker/internal/constants"
	"github.com/Sinanaas/gotth-financial-tracker/internal/controllers"
	"github.com/Sinanaas/gotth-financial-tracker/internal/templates"
	"github.com/Sinanaas/gotth-financial-tracker/internal/utils"
	"github.com/gin-gonic/gin"
	"net/http"
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

	categories, err := h.BC.GetAllCategories()
	if err != nil {
		return
	}

	transactions, err := h.BC.GetUserTransactions(userId)
	if err != nil {
		return
	}

	accounts, err := h.BC.GetUserAccounts(userId)
	if err != nil {
		return
	}

	var transactionType constants.TransactionType
	transactionTypeArray := transactionType.ToArrayString()

	c := templates.Transaction(categories, transactions, accounts, transactionTypeArray)
	err = templates.Layout(c, cookie).Render(ctx.Request.Context(), ctx.Writer)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
}
