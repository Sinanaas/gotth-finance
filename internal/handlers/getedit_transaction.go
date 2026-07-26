package handlers

import (
	"net/http"

	"github.com/Sinanaas/gotth-financial-tracker/internal/managers"
	"github.com/Sinanaas/gotth-financial-tracker/internal/templates"
	"github.com/Sinanaas/gotth-financial-tracker/internal/utils"
	"github.com/gin-gonic/gin"
)

type GetEditTransactionHandler struct {
	BM *managers.BasicManager
}

func NewGetEditTransactionHandler(bm *managers.BasicManager) *GetEditTransactionHandler {
	return &GetEditTransactionHandler{BM: bm}
}

func (h *GetEditTransactionHandler) ServeHTTP(ctx *gin.Context) {
	id := ctx.Param("id")
	userId := utils.GetSessionUserID(ctx)

	transaction, err := h.BM.FindTransactionById(id)
	if err != nil {
		swalError(ctx, "Transaction not found")
		return
	}
	if transaction.UserID.String() != userId {
		ctx.Status(http.StatusForbidden)
		return
	}

	categories, err := h.BM.GetUserCategories(userId)
	if err != nil {
		swalError(ctx, "Failed to load categories")
		return
	}

	accounts, err := h.BM.GetUserAccounts(userId)
	if err != nil {
		swalError(ctx, "Failed to load accounts")
		return
	}

	c := templates.EditTransactionModal(transaction, categories, accounts)
	if err := c.Render(ctx.Request.Context(), ctx.Writer); err != nil {
		swalError(ctx, "Failed to render modal")
	}
}
