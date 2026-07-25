package handlers

import (
	"github.com/Sinanaas/gotth-financial-tracker/internal/managers"
	"github.com/Sinanaas/gotth-financial-tracker/internal/utils"
	"github.com/Sinanaas/gotth-financial-tracker/internal/validation"
	"github.com/gin-gonic/gin"
)

type PostGoalContributeHandler struct{ BM *managers.BasicManager }

func NewPostGoalContributeHandler(bm *managers.BasicManager) *PostGoalContributeHandler {
	return &PostGoalContributeHandler{BM: bm}
}

func (h *PostGoalContributeHandler) ServeHTTP(ctx *gin.Context) {
	userId := utils.GetSessionUserID(ctx)
	id, err := validation.UUIDField(ctx.Param("id"), "goal")
	if err != nil {
		swalError(ctx, err.Error())
		return
	}
	amount, err := validation.Amount(ctx.PostForm("Amount"))
	if err != nil {
		swalError(ctx, err.Error())
		return
	}
	if err := h.BM.AddToGoal(id, userId, amount); err != nil {
		swalError(ctx, "Failed to add contribution")
		return
	}
	renderGoalsPanel(ctx, h.BM, userId)
}
