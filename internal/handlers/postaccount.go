package handlers

import (
	"strconv"

	"github.com/Sinanaas/gotth-financial-tracker/internal/managers"
	"github.com/Sinanaas/gotth-financial-tracker/internal/models"
	"github.com/Sinanaas/gotth-financial-tracker/internal/templates"
	"github.com/Sinanaas/gotth-financial-tracker/internal/utils"
	"github.com/Sinanaas/gotth-financial-tracker/internal/validation"
	"github.com/gin-gonic/gin"
)

type PostAccountHandler struct {
	BM *managers.BasicManager
}

func NewPostAccountHandler(bm *managers.BasicManager) *PostAccountHandler {
	return &PostAccountHandler{BM: bm}
}

func (h *PostAccountHandler) ServeHTTP(ctx *gin.Context) {
	userId := utils.GetSessionUserID(ctx)

	name, err := validation.Required(ctx.PostForm("Name"), "account name")
	if err != nil {
		swalError(ctx, err.Error())
		return
	}
	balance, err := strconv.ParseFloat(ctx.PostForm("Balance"), 64)
	if err != nil || balance < 0 {
		swalError(ctx, "balance must be zero or more")
		return
	}

	payload := models.AccountRequest{
		Name:        name,
		Description: ctx.PostForm("Description"),
		Balance:     balance,
		UserID:      userId,
	}
	if err := h.BM.CreateAccount(payload); err != nil {
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
