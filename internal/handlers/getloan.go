package handlers

import (
	"github.com/Sinanaas/gotth-financial-tracker/internal/constants"
	"github.com/Sinanaas/gotth-financial-tracker/internal/controllers"
	"github.com/Sinanaas/gotth-financial-tracker/internal/templates"
	"github.com/Sinanaas/gotth-financial-tracker/internal/utils"
	"github.com/gin-gonic/gin"
	"net/http"
)

type GetLoanHandler struct {
	BC *controllers.BasicController
}

func NewGetLoanHandler(bc *controllers.BasicController) *GetLoanHandler {
	return &GetLoanHandler{BC: bc}
}

func (h *GetLoanHandler) ServeHTTP(ctx *gin.Context) {
	cookie, _ := ctx.Cookie("access_token")
	userId := utils.GetSessionUserID(ctx)

	categories, err := h.BC.GetAllCategories()
	if err != nil {
		return
	}

	loans, err := h.BC.GetLoanWithCategoryName(userId)
	if err != nil {
		return
	}

	accounts, err := h.BC.GetUserAccounts(userId)
	if err != nil {
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
