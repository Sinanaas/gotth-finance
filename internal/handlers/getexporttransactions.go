package handlers

import (
	"encoding/csv"
	"net/http"
	"strconv"

	"github.com/Sinanaas/gotth-financial-tracker/internal/managers"
	"github.com/Sinanaas/gotth-financial-tracker/internal/utils"
	"github.com/gin-gonic/gin"
)

type ExportTransactionsHandler struct {
	BM *managers.BasicManager
}

func NewExportTransactionsHandler(bm *managers.BasicManager) *ExportTransactionsHandler {
	return &ExportTransactionsHandler{BM: bm}
}

// CSV column contract (also the import format for round-tripping):
// date,type,description,category,account,amount
func (h *ExportTransactionsHandler) ServeHTTP(ctx *gin.Context) {
	userId := utils.GetSessionUserID(ctx)

	transactions, err := h.BM.GetUserTransactions(userId)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to load transactions"})
		return
	}

	ctx.Header("Content-Type", "text/csv")
	ctx.Header("Content-Disposition", `attachment; filename="transactions.csv"`)

	w := csv.NewWriter(ctx.Writer)
	_ = w.Write([]string{"date", "type", "description", "category", "account", "amount"})
	for _, t := range transactions {
		_ = w.Write([]string{
			t.TransactionDate.Format("2006-01-02"),
			t.TransactionType.ToString(),
			t.Description,
			t.Category.Name,
			t.Account.Name,
			strconv.FormatFloat(t.Amount, 'f', 2, 64),
		})
	}
	w.Flush()
}
