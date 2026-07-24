package handlers

import (
	"github.com/Sinanaas/gotth-financial-tracker/internal/managers"
	"github.com/gin-gonic/gin"
)

type DeleteRecurringHandler struct {
	BM *managers.BasicManager
}

func NewDeleteRecurringHandler(bm *managers.BasicManager) *DeleteRecurringHandler {
	return &DeleteRecurringHandler{BM: bm}
}

func (h *DeleteRecurringHandler) ServeHTTP(ctx *gin.Context) {
	recurringId := ctx.PostForm("RecurringID")
	if recurringId == "" {
		ctx.JSON(400, gin.H{"error": "Recurring ID is missing"})
		return
	}

	err := h.BM.DeleteRecurringById(recurringId)
	if err != nil {
		ctx.JSON(500, gin.H{"error": "Failed to delete recurring"})
		return
	}
	ctx.Header("HX-Redirect", "/recurring")
}
