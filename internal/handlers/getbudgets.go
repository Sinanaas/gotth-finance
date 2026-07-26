package handlers

import (
	"github.com/Sinanaas/gotth-financial-tracker/internal/managers"
	"github.com/Sinanaas/gotth-financial-tracker/internal/templates"
	"github.com/Sinanaas/gotth-financial-tracker/internal/utils"
	"github.com/gin-gonic/gin"
)

type GetBudgetsHandler struct {
	BM *managers.BasicManager
}

func NewGetBudgetsHandler(bm *managers.BasicManager) *GetBudgetsHandler {
	return &GetBudgetsHandler{BM: bm}
}

func (h *GetBudgetsHandler) ServeHTTP(ctx *gin.Context) {
	cookie, _ := ctx.Cookie("access_token")
	userId := utils.GetSessionUserID(ctx)

	statuses, err := h.BM.GetBudgetStatus(userId)
	if err != nil {
		renderErrorPage(ctx, cookie, "Failed to load budgets")
		return
	}
	categories, err := h.BM.GetUserCategories(userId)
	if err != nil {
		renderErrorPage(ctx, cookie, "Failed to load categories")
		return
	}

	c := templates.Budgets(statuses, categories)
	if err := templates.Layout(c, cookie).Render(ctx.Request.Context(), ctx.Writer); err != nil {
		renderErrorPage(ctx, cookie, "Failed to render page")
	}
}
