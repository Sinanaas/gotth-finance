package handlers

import (
	"net/http"

	"github.com/Sinanaas/gotth-financial-tracker/internal/managers"
	"github.com/Sinanaas/gotth-financial-tracker/internal/templates"
	"github.com/Sinanaas/gotth-financial-tracker/internal/utils"
	"github.com/gin-gonic/gin"
)

type GetEditAccountHandler struct {
	BM *managers.BasicManager
}

func NewGetEditAccountHandler(bm *managers.BasicManager) *GetEditAccountHandler {
	return &GetEditAccountHandler{BM: bm}
}

func (h *GetEditAccountHandler) ServeHTTP(ctx *gin.Context) {
	userId := utils.GetSessionUserID(ctx)
	id := ctx.Param("id")

	account, err := h.BM.FindUserAccountById(userId, id)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Account not found"})
		return
	}

	if err := templates.EditAccountModal(account).Render(ctx.Request.Context(), ctx.Writer); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to render modal"})
	}
}
