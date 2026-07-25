package handlers

import (
	"net/http"
	"strconv"

	"github.com/Sinanaas/gotth-financial-tracker/internal/managers"
	"github.com/Sinanaas/gotth-financial-tracker/internal/templates"
	"github.com/Sinanaas/gotth-financial-tracker/internal/utils"
	"github.com/gin-gonic/gin"
)

const pageSize = 20

type GetTransactionListHandler struct {
	BM *managers.BasicManager
}

func NewGetTransactionListHandler(bm *managers.BasicManager) *GetTransactionListHandler {
	return &GetTransactionListHandler{BM: bm}
}

func (h *GetTransactionListHandler) ServeHTTP(ctx *gin.Context) {
	userId := utils.GetSessionUserID(ctx)

	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	if page < 1 {
		page = 1
	}

	startDate := ctx.Query("StartDate")
	endDate := ctx.Query("EndDate")
	categoryID := ctx.Query("Category")
	accountID := ctx.Query("Account")
	search := ctx.Query("Search")

	txTypeStr := ctx.DefaultQuery("Type", "-1")
	txType, err := strconv.Atoi(txTypeStr)
	if err != nil {
		txType = -1
	}

	transactions, total, err := h.BM.FilterTransactions(userId, startDate, endDate, categoryID, accountID, search, txType, page, pageSize)
	if err != nil {
		swalError(ctx, "Failed to load transactions")
		return
	}

	c := templates.TransactionListBody(transactions, total, page, pageSize)
	if err := c.Render(ctx.Request.Context(), ctx.Writer); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to render"})
	}
}
