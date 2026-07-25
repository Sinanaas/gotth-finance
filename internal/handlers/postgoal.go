package handlers

import (
	"github.com/Sinanaas/gotth-financial-tracker/internal/managers"
	"github.com/Sinanaas/gotth-financial-tracker/internal/models"
	"github.com/Sinanaas/gotth-financial-tracker/internal/templates"
	"github.com/Sinanaas/gotth-financial-tracker/internal/utils"
	"github.com/Sinanaas/gotth-financial-tracker/internal/validation"
	"github.com/gin-gonic/gin"
)

type PostGoalHandler struct{ BM *managers.BasicManager }

func NewPostGoalHandler(bm *managers.BasicManager) *PostGoalHandler { return &PostGoalHandler{BM: bm} }

func (h *PostGoalHandler) ServeHTTP(ctx *gin.Context) {
	userId := utils.GetSessionUserID(ctx)
	name, err := validation.Required(ctx.PostForm("Name"), "name")
	if err != nil {
		swalError(ctx, err.Error())
		return
	}
	target, err := validation.Amount(ctx.PostForm("Target"))
	if err != nil {
		swalError(ctx, err.Error())
		return
	}
	accountId, err := validation.UUIDField(ctx.PostForm("Account"), "account")
	if err != nil {
		swalError(ctx, err.Error())
		return
	}
	if err := h.BM.CreateGoal(models.GoalRequest{Name: name, TargetAmount: target, AccountID: accountId, UserID: userId}); err != nil {
		swalError(ctx, err.Error())
		return
	}
	renderGoalsPanel(ctx, h.BM, userId)
}

func renderGoalsPanel(ctx *gin.Context, bm *managers.BasicManager, userId string) {
	statuses, err := bm.GetGoalStatus(userId)
	if err != nil {
		swalError(ctx, "Failed to reload goals")
		return
	}
	accounts, err := bm.GetUserAccounts(userId)
	if err != nil {
		swalError(ctx, "Failed to reload accounts")
		return
	}
	ctx.Status(200)
	_ = templates.GoalsPanel(statuses, accounts).Render(ctx.Request.Context(), ctx.Writer)
}
