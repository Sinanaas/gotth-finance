package handlers

import (
	"github.com/Sinanaas/gotth-financial-tracker/internal/managers"
	"github.com/Sinanaas/gotth-financial-tracker/internal/utils"
	"github.com/Sinanaas/gotth-financial-tracker/internal/validation"
	"github.com/gin-gonic/gin"
)

type DeleteGoalHandler struct{ BM *managers.BasicManager }

func NewDeleteGoalHandler(bm *managers.BasicManager) *DeleteGoalHandler { return &DeleteGoalHandler{BM: bm} }

func (h *DeleteGoalHandler) ServeHTTP(ctx *gin.Context) {
	userId := utils.GetSessionUserID(ctx)
	id, err := validation.UUIDField(ctx.Param("id"), "goal")
	if err != nil {
		swalError(ctx, err.Error())
		return
	}
	if err := h.BM.DeleteGoal(id, userId); err != nil {
		swalError(ctx, "Failed to delete goal")
		return
	}
	renderGoalsPanel(ctx, h.BM, userId)
}
