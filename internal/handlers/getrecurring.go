package handlers

import (
	"net/http"

	"github.com/Sinanaas/gotth-financial-tracker/internal/constants"
	"github.com/Sinanaas/gotth-financial-tracker/internal/managers"
	"github.com/Sinanaas/gotth-financial-tracker/internal/templates"
	"github.com/Sinanaas/gotth-financial-tracker/internal/utils"
	"github.com/gin-gonic/gin"
)

type GetRecurringHandler struct {
	BM *managers.BasicManager
}

func NewGetRecurringHandler(bm *managers.BasicManager) *GetRecurringHandler {
	return &GetRecurringHandler{BM: bm}
}

func (h *GetRecurringHandler) ServeHTTP(ctx *gin.Context) {
	cookie, _ := ctx.Cookie("access_token")
	userId := utils.GetSessionUserID(ctx)

	categories, err := h.BM.GetUserCategories(userId)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to load categories"})
		return
	}

	recurring, err := h.BM.GetRecurrings(userId)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to load recurring"})
		return
	}

	accounts, err := h.BM.GetUserAccounts(userId)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to load accounts"})
		return
	}

	var transactionType constants.TransactionType
	transactionTypeArray := transactionType.ToArrayString()

	var periodicity constants.Periodicity
	periodicityArray := periodicity.ToArrayString()

	c := templates.Recurring(categories, recurring, accounts, transactionTypeArray, periodicityArray)
	err = templates.Layout(c, cookie).Render(ctx.Request.Context(), ctx.Writer)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
}
