package handlers

import (
	"net/http"

	"github.com/Sinanaas/gotth-financial-tracker/internal/constants"
	"github.com/Sinanaas/gotth-financial-tracker/internal/managers"
	"github.com/Sinanaas/gotth-financial-tracker/internal/templates"
	"github.com/Sinanaas/gotth-financial-tracker/internal/utils"
	"github.com/gin-gonic/gin"
)

type GetEditRecurringHandler struct {
	BM *managers.BasicManager
}

func NewGetEditRecurringHandler(bm *managers.BasicManager) *GetEditRecurringHandler {
	return &GetEditRecurringHandler{BM: bm}
}

func (h *GetEditRecurringHandler) ServeHTTP(ctx *gin.Context) {
	userId := utils.GetSessionUserID(ctx)
	id := ctx.Param("id")

	recurring, err := h.BM.FindUserRecurringById(userId, id)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Recurring not found"})
		return
	}
	categories, err := h.BM.GetUserCategories(userId)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to load categories"})
		return
	}
	accounts, err := h.BM.GetUserAccounts(userId)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to load accounts"})
		return
	}

	var txType constants.TransactionType
	var periodicity constants.Periodicity

	c := templates.EditRecurringModal(recurring, categories, accounts, txType.ToArrayString(), periodicity.ToArrayString())
	if err := c.Render(ctx.Request.Context(), ctx.Writer); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to render modal"})
	}
}
