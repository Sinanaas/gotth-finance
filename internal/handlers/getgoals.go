package handlers

import (
	"github.com/Sinanaas/gotth-financial-tracker/internal/managers"
	"github.com/Sinanaas/gotth-financial-tracker/internal/templates"
	"github.com/Sinanaas/gotth-financial-tracker/internal/utils"
	"github.com/gin-gonic/gin"
)

type GetGoalsHandler struct{ BM *managers.BasicManager }

func NewGetGoalsHandler(bm *managers.BasicManager) *GetGoalsHandler { return &GetGoalsHandler{BM: bm} }

func (h *GetGoalsHandler) ServeHTTP(ctx *gin.Context) {
	cookie, _ := ctx.Cookie("access_token")
	userId := utils.GetSessionUserID(ctx)
	statuses, err := h.BM.GetGoalStatus(userId)
	if err != nil {
		renderErrorPage(ctx, cookie, "Failed to load goals")
		return
	}
	accounts, err := h.BM.GetUserAccounts(userId)
	if err != nil {
		renderErrorPage(ctx, cookie, "Failed to load accounts")
		return
	}
	if err := templates.Layout(templates.Goals(statuses, accounts), cookie).Render(ctx.Request.Context(), ctx.Writer); err != nil {
		renderErrorPage(ctx, cookie, "Failed to render page")
	}
}
