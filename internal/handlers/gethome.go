package handlers

import (
	"net/http"
	"time"

	"github.com/Sinanaas/gotth-financial-tracker/internal/managers"
	"github.com/Sinanaas/gotth-financial-tracker/internal/templates"
	"github.com/Sinanaas/gotth-financial-tracker/internal/utils"
	"github.com/gin-gonic/gin"
)

type GetHomeHandler struct {
	BM *managers.BasicManager
}

func NewGetHomeHandler(bm *managers.BasicManager) *GetHomeHandler {
	return &GetHomeHandler{BM: bm}
}

func (h *GetHomeHandler) ServeHTTP(ctx *gin.Context) {
	userId := utils.GetSessionUserID(ctx)
	cookie, _ := ctx.Cookie("access_token")
	now := time.Now()

	// Current month totals
	thisMonthIncome, thisMonthExpenses, err := h.BM.GetUserMonthlyTotals(userId, now.Year(), now.Month())
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to load monthly totals"})
		return
	}

	// Previous month totals for MoM delta
	prev := now.AddDate(0, -1, 0)
	prevIncome, prevExpense, err := h.BM.GetUserMonthlyTotals(userId, prev.Year(), prev.Month())
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to load previous month totals"})
		return
	}

	accounts, err := h.BM.GetUserAccounts(userId)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to load accounts"})
		return
	}

	loans, err := h.BM.GetUserActiveLoans(userId)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to load loans"})
		return
	}

	totalBalance := h.BM.GetUserTotalBalance(userId)

	transactions, err := h.BM.GetUserLatestSixTransactions(userId)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to load transactions"})
		return
	}

	// All upcoming recurring (not just closest)
	recurrings, err := h.BM.GetAllUpcomingRecurring(userId)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to load recurring"})
		return
	}

	topCategories, err := h.BM.GetUserTopCategories(userId)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to load top categories"})
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
		inc, exp, _ := h.BM.GetUserMonthlyTotals(userId, t.Year(), t.Month())
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
