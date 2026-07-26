package handlers

import (
	"github.com/Sinanaas/gotth-financial-tracker/internal/managers"
	"github.com/Sinanaas/gotth-financial-tracker/internal/models"
	"github.com/Sinanaas/gotth-financial-tracker/internal/templates"
	"github.com/Sinanaas/gotth-financial-tracker/internal/utils"
	"github.com/Sinanaas/gotth-financial-tracker/internal/validation"
	"github.com/gin-gonic/gin"
)

type PostBudgetHandler struct {
	BM *managers.BasicManager
}

func NewPostBudgetHandler(bm *managers.BasicManager) *PostBudgetHandler {
	return &PostBudgetHandler{BM: bm}
}

func (h *PostBudgetHandler) ServeHTTP(ctx *gin.Context) {
	userId := utils.GetSessionUserID(ctx)

	categoryId, err := validation.UUIDField(ctx.PostForm("Category"), "category")
	if err != nil {
		swalError(ctx, err.Error())
		return
	}
	amount, err := validation.Amount(ctx.PostForm("Amount"))
	if err != nil {
		swalError(ctx, err.Error())
		return
	}

	if err := h.BM.CreateBudget(models.BudgetRequest{CategoryID: categoryId, Amount: amount, UserID: userId}); err != nil {
		swalError(ctx, err.Error())
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
