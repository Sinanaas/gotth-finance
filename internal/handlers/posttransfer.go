package handlers

import (
	"github.com/Sinanaas/gotth-financial-tracker/internal/managers"
	"github.com/Sinanaas/gotth-financial-tracker/internal/models"
	"github.com/Sinanaas/gotth-financial-tracker/internal/templates"
	"github.com/Sinanaas/gotth-financial-tracker/internal/utils"
	"github.com/Sinanaas/gotth-financial-tracker/internal/validation"
	"github.com/gin-gonic/gin"
)

type PostTransferHandler struct {
	BM *managers.BasicManager
}

func NewPostTransferHandler(bm *managers.BasicManager) *PostTransferHandler {
	return &PostTransferHandler{BM: bm}
}

func (h *PostTransferHandler) ServeHTTP(ctx *gin.Context) {
	userId := utils.GetSessionUserID(ctx)

	fromId, err := validation.UUIDField(ctx.PostForm("FromAccount"), "source account")
	if err != nil {
		swalError(ctx, err.Error())
		return
	}
	toId, err := validation.UUIDField(ctx.PostForm("ToAccount"), "destination account")
	if err != nil {
		swalError(ctx, err.Error())
		return
	}
	if fromId == toId {
		swalError(ctx, "cannot transfer to the same account")
		return
	}
	amount, err := validation.Amount(ctx.PostForm("Amount"))
	if err != nil {
		swalError(ctx, err.Error())
		return
	}
	if _, err := validation.Date(ctx.PostForm("Date")); err != nil {
		swalError(ctx, err.Error())
		return
	}

	payload := models.TransferRequest{
		FromAccountID: fromId,
		ToAccountID:   toId,
		Amount:        amount,
		Date:          ctx.PostForm("Date"),
		Description:   ctx.PostForm("Description"),
		UserID:        userId,
	}
	if err := h.BM.CreateTransfer(payload); err != nil {
		swalError(ctx, err.Error())
		return
	}

	accounts, err := h.BM.GetUserAccounts(userId)
	if err != nil {
		swalError(ctx, "Failed to reload accounts")
		return
	}
	ctx.Status(200)
	_ = templates.AccountsPanel(accounts).Render(ctx.Request.Context(), ctx.Writer)
}
