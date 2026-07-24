package handlers

import (
	"time"

	"github.com/Sinanaas/gotth-financial-tracker/internal/controllers"
	"github.com/Sinanaas/gotth-financial-tracker/internal/templates"
	"github.com/Sinanaas/gotth-financial-tracker/internal/utils"
	"github.com/gin-gonic/gin"
)

type GetHomeHandler struct {
	BC *controllers.BasicController
}

func NewGetHomeHandler(bc *controllers.BasicController) *GetHomeHandler {
	return &GetHomeHandler{BC: bc}
}

func (h *GetHomeHandler) ServeHTTP(ctx *gin.Context) {
	userId := utils.GetSessionUserID(ctx)
	cookie, _ := ctx.Cookie("access_token")
	now := time.Now()

	// Current month totals
	thisMonthIncome, thisMonthExpenses, err := h.BC.GetUserMonthlyTotals(userId, now.Year(), now.Month())
	if err != nil {
		return
	}

	// Previous month totals for MoM delta
	prev := now.AddDate(0, -1, 0)
	prevIncome, prevExpense, err := h.BC.GetUserMonthlyTotals(userId, prev.Year(), prev.Month())
	if err != nil {
		return
	}

	accounts, err := h.BC.GetUserAccounts(userId)
	if err != nil {
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

	// All upcoming recurring (not just closest)
	recurrings, err := h.BC.GetAllUpcomingRecurring(userId)
	if err != nil {
		return
	}

	topCategories, err := h.BC.GetUserTopCategories(userId)
	if err != nil {
		return
	}

	// Build 6-month chart data
	monthlyLabels := make([]string, 6)
	monthlyIncome := make([]float64, 6)
	monthlyExpense := make([]float64, 6)
	for i := 5; i >= 0; i-- {
		t := now.AddDate(0, -i, 0)
		idx := 5 - i
		monthlyLabels[idx] = t.Format("Jan '06")
		inc, exp, _ := h.BC.GetUserMonthlyTotals(userId, t.Year(), t.Month())
		monthlyIncome[idx] = inc
		monthlyExpense[idx] = exp
	}

	c := templates.Home(
		thisMonthIncome,
		thisMonthExpenses,
		prevIncome,
		prevExpense,
		accounts,
		loans,
		transactions,
		recurrings,
		now.Month().String(),
		utils.FormatCurrency(totalBalance),
		topCategories,
		monthlyLabels,
		monthlyIncome,
		monthlyExpense,
	)
	if err := templates.Layout(c, cookie).Render(ctx.Request.Context(), ctx.Writer); err != nil {
		ctx.JSON(500, gin.H{"error": "Error rendering template"})
	}
}
