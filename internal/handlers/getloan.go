package handlers

import (
	"net/http"

	"github.com/Sinanaas/gotth-financial-tracker/internal/constants"
	"github.com/Sinanaas/gotth-financial-tracker/internal/managers"
	"github.com/Sinanaas/gotth-financial-tracker/internal/templates"
	"github.com/Sinanaas/gotth-financial-tracker/internal/utils"
	"github.com/gin-gonic/gin"
)

type GetLoanHandler struct {
	BM *managers.BasicManager
}

func NewGetLoanHandler(bm *managers.BasicManager) *GetLoanHandler {
	return &GetLoanHandler{BM: bm}
}

func (h *GetLoanHandler) ServeHTTP(ctx *gin.Context) {
	cookie, _ := ctx.Cookie("access_token")
	userId := utils.GetSessionUserID(ctx)

	categories, err := h.BM.GetUserCategories(userId)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to load categories"})
		return
	}

	loans, err := h.BM.GetLoans(userId)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to load loans"})
		return
	}

	accounts, err := h.BM.GetUserAccounts(userId)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to load accounts"})
		return
	}

	var transactionType constants.TransactionType
	transactionTypeArray := transactionType.ToArrayString()

	c := templates.Loans(categories, loans, accounts, transactionTypeArray)

	err = templates.Layout(c, cookie).Render(ctx.Request.Context(), ctx.Writer)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
}
