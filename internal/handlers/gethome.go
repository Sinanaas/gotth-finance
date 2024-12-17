package handlers

import (
	"github.com/Sinanaas/gotth-financial-tracker/internal/controllers"
	"github.com/Sinanaas/gotth-financial-tracker/internal/templates"
	"github.com/Sinanaas/gotth-financial-tracker/internal/utils"
	"github.com/gin-gonic/gin"
	"time"
)

type GetHomeHandler struct {
	BC *controllers.BasicController
}

func NewGetHomeHandler(bc *controllers.BasicController) *GetHomeHandler {
	return &GetHomeHandler{BC: bc}
}

func (h *GetHomeHandler) ServeHTTP(ctx *gin.Context) {
	userId := utils.GetSessionUserID(ctx)

	accounts, err := h.BC.GetUserAccounts(userId)
	if err != nil {
		return
	}

	thisMonthExpenses, err := h.BC.GetUserMonthlyExpenses(userId)
	if err != nil && thisMonthExpenses > 0 {
		return
	}

	thisMonthIncome, err := h.BC.GetUserMonthlyIncome(userId)
	if err != nil && thisMonthIncome > 0 {
		return
	}

	loans, err := h.BC.GetUserActiveLoans(userId)
	if err != nil {
		return
	}

	totalBalance := h.BC.GetUserTotalBalance(userId)

	transactions, err := h.BC.GetUserLatestSixTransactions(userId)
	if err != nil {
		return
	}

	recurring, err := h.BC.GetUserUpcomingRecurring(userId)
	if err != nil {
		return
	}

	topCategories, err := h.BC.GetUserTopCategories(userId)
	if err != nil {
		return
	}

	cookie, _ := ctx.Cookie("access_token")
	month := time.Now().Month()
	c := templates.Home(thisMonthIncome, thisMonthExpenses, accounts, loans, transactions, recurring, month.String(), utils.FormatCurrency(totalBalance), topCategories)
	err = templates.Layout(c, cookie).Render(ctx.Request.Context(), ctx.Writer)
	if err != nil {
		ctx.JSON(500, gin.H{"error": "Error rendering template"})
		return
	}
}
