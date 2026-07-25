package handlers

import (
	"github.com/Sinanaas/gotth-financial-tracker/internal/managers"
	"github.com/Sinanaas/gotth-financial-tracker/internal/templates"
	"github.com/Sinanaas/gotth-financial-tracker/internal/utils"
	"github.com/Sinanaas/gotth-financial-tracker/internal/validation"
	"github.com/gin-gonic/gin"
)

type DeleteBudgetHandler struct {
	BM *managers.BasicManager
}

func NewDeleteBudgetHandler(bm *managers.BasicManager) *DeleteBudgetHandler {
	return &DeleteBudgetHandler{BM: bm}
}

func (h *DeleteBudgetHandler) ServeHTTP(ctx *gin.Context) {
	userId := utils.GetSessionUserID(ctx)

	budgetId, err := validation.UUIDField(ctx.Param("id"), "budget")
	if err != nil {
		swalError(ctx, err.Error())
		return
	}

	if err := h.BM.DeleteBudget(budgetId, userId); err != nil {
		swalError(ctx, "Failed to delete budget")
		return
	}

	statuses, err := h.BM.GetBudgetStatus(userId)
	if err != nil {
		swalError(ctx, "Failed to reload budgets")
		return
	}
	ctx.Status(200)
	_ = templates.BudgetsPanel(statuses).Render(ctx.Request.Context(), ctx.Writer)
}
