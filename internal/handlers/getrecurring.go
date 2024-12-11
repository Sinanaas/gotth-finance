package handlers

import (
	"github.com/Sinanaas/gotth-financial-tracker/internal/constants"
	"github.com/Sinanaas/gotth-financial-tracker/internal/controllers"
	"github.com/Sinanaas/gotth-financial-tracker/internal/templates"
	"github.com/Sinanaas/gotth-financial-tracker/internal/utils"
	"github.com/gin-gonic/gin"
	"net/http"
)

type GetRecurringHandler struct {
	BC *controllers.BasicController
}

func NewGetRecurringHandler(bc *controllers.BasicController) *GetRecurringHandler {
	return &GetRecurringHandler{BC: bc}
}

func (h *GetRecurringHandler) ServeHTTP(ctx *gin.Context) {
	cookie, _ := ctx.Cookie("access_token")
	userId := utils.GetSessionUserID(ctx)

	categories, err := h.BC.GetAllCategories()
	if err != nil {
		return
	}

	recurring, err := h.BC.GetRecurringWithCategoryName(userId)
	if err != nil {
		return
	}

	accounts, err := h.BC.GetUserAccounts(userId)
	if err != nil {
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
