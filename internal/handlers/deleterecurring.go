package handlers

import (
	"github.com/Sinanaas/gotth-financial-tracker/internal/managers"
	"github.com/Sinanaas/gotth-financial-tracker/internal/templates"
	"github.com/Sinanaas/gotth-financial-tracker/internal/utils"
	"github.com/Sinanaas/gotth-financial-tracker/internal/validation"
	"github.com/gin-gonic/gin"
)

type DeleteRecurringHandler struct {
	BM *managers.BasicManager
}

func NewDeleteRecurringHandler(bm *managers.BasicManager) *DeleteRecurringHandler {
	return &DeleteRecurringHandler{BM: bm}
}

func (h *DeleteRecurringHandler) ServeHTTP(ctx *gin.Context) {
	userId := utils.GetSessionUserID(ctx)

	recurringId, err := validation.UUIDField(ctx.PostForm("RecurringID"), "recurring payment")
	if err != nil {
		swalError(ctx, err.Error())
		return
	}

	if err := h.BM.DeleteRecurringById(recurringId, userId); err != nil {
		swalError(ctx, "Failed to delete recurring payment")
		return
	}

	recurrings, err := h.BM.GetRecurrings(userId)
	if err != nil {
		swalError(ctx, "Failed to reload recurring")
		return
	}
	ctx.Status(200)
	_ = templates.RecurringPanel(recurrings).Render(ctx.Request.Context(), ctx.Writer)
}
