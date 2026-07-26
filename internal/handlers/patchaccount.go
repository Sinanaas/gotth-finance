package handlers

import (
	"github.com/Sinanaas/gotth-financial-tracker/internal/managers"
	"github.com/Sinanaas/gotth-financial-tracker/internal/templates"
	"github.com/Sinanaas/gotth-financial-tracker/internal/utils"
	"github.com/Sinanaas/gotth-financial-tracker/internal/validation"
	"github.com/gin-gonic/gin"
)

type PatchAccountHandler struct {
	BM *managers.BasicManager
}

func NewPatchAccountHandler(bm *managers.BasicManager) *PatchAccountHandler {
	return &PatchAccountHandler{BM: bm}
}

func (h *PatchAccountHandler) ServeHTTP(ctx *gin.Context) {
	userId := utils.GetSessionUserID(ctx)
	id := ctx.Param("id")

	name, err := validation.Required(ctx.PostForm("Name"), "account name")
	if err != nil {
		swalError(ctx, err.Error())
		return
	}

	if err := h.BM.UpdateAccount(userId, id, name, ctx.PostForm("Description")); err != nil {
		swalError(ctx, "Failed to update account")
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
