package handlers

import (
	"github.com/Sinanaas/gotth-financial-tracker/internal/managers"
	"github.com/Sinanaas/gotth-financial-tracker/internal/templates"
	"github.com/Sinanaas/gotth-financial-tracker/internal/utils"
	"github.com/Sinanaas/gotth-financial-tracker/internal/validation"
	"github.com/gin-gonic/gin"
)

type DeleteAccountHandler struct {
	BM *managers.BasicManager
}

func NewDeleteAccountHandler(bm *managers.BasicManager) *DeleteAccountHandler {
	return &DeleteAccountHandler{BM: bm}
}

func (h *DeleteAccountHandler) ServeHTTP(ctx *gin.Context) {
	userId := utils.GetSessionUserID(ctx)

	accountId, err := validation.UUIDField(ctx.PostForm("AccountID"), "account")
	if err != nil {
		swalError(ctx, err.Error())
		return
	}

	if err := h.BM.DeleteAccountById(accountId, userId); err != nil {
		swalError(ctx, "Failed to delete account")
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
