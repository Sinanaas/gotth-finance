package handlers

import (
	"github.com/Sinanaas/gotth-financial-tracker/internal/managers"
	"github.com/Sinanaas/gotth-financial-tracker/internal/models"
	"github.com/Sinanaas/gotth-financial-tracker/internal/templates"
	"github.com/Sinanaas/gotth-financial-tracker/internal/utils"
	"github.com/Sinanaas/gotth-financial-tracker/internal/validation"
	"github.com/gin-gonic/gin"
)

type PostRecurringHandler struct {
	BM *managers.BasicManager
}

func NewPostRecurringHandler(bm *managers.BasicManager) *PostRecurringHandler {
	return &PostRecurringHandler{BM: bm}
}

func (h *PostRecurringHandler) ServeHTTP(ctx *gin.Context) {
	userId := utils.GetSessionUserID(ctx)

	name, err := validation.Required(ctx.PostForm("Name"), "name")
	if err != nil {
		swalError(ctx, err.Error())
		return
	}
	amount, err := validation.Amount(ctx.PostForm("Amount"))
	if err != nil {
		swalError(ctx, err.Error())
		return
	}
	txType, err := validation.IntField(ctx.PostForm("Type"), "type")
	if err != nil {
		swalError(ctx, err.Error())
		return
	}
	periodicity, err := validation.IntField(ctx.PostForm("Periodicity"), "periodicity")
	if err != nil {
		swalError(ctx, err.Error())
		return
	}
	categoryId, err := validation.UUIDField(ctx.PostForm("Category"), "category")
	if err != nil {
		swalError(ctx, err.Error())
		return
	}
	accountId, err := validation.UUIDField(ctx.PostForm("Account"), "account")
	if err != nil {
		swalError(ctx, err.Error())
		return
	}
	if _, err := validation.Date(ctx.PostForm("StartDate")); err != nil {
		swalError(ctx, err.Error())
		return
	}

	payload := models.RecurringRequest{
		Name:            name,
		Amount:          amount,
		TransactionType: txType,
		Periodicity:     periodicity,
		CategoryID:      categoryId,
		AccountID:       accountId,
		StartDate:       ctx.PostForm("StartDate"),
		UserID:          userId,
	}
	if err := h.BM.CreateRecurring(payload); err != nil {
		swalError(ctx, err.Error())
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
