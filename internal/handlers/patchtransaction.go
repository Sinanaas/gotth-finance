package handlers

import (
	"net/http"
	"strconv"

	"github.com/Sinanaas/gotth-financial-tracker/internal/managers"
	"github.com/Sinanaas/gotth-financial-tracker/internal/models"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

type PatchTransactionHandler struct {
	BM *managers.BasicManager
}

func NewPatchTransactionHandler(bm *managers.BasicManager) *PatchTransactionHandler {
	return &PatchTransactionHandler{BM: bm}
}

func (h *PatchTransactionHandler) ServeHTTP(ctx *gin.Context) {
	var payload models.TransactionUpdateRequest
	var err error

	payload.ID = ctx.Param("id")
	payload.Description = ctx.PostForm("Description")
	payload.CategoryID = ctx.PostForm("Category")
	payload.OldAccountID = ctx.PostForm("OldAccountID")
	payload.NewAccountID = ctx.PostForm("Account")
	payload.Amount, err = strconv.ParseFloat(ctx.PostForm("Amount"), 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	payload.Date = ctx.PostForm("Date")
	payload.Type, err = strconv.Atoi(ctx.PostForm("Type"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	session := sessions.Default(ctx)
	var userId string
	v := session.Get("user_id")
	if v != nil {
		userId = v.(string)
	}
	payload.UserID = userId

	if err := h.BM.UpdateTransaction(payload); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.Header("HX-Redirect", "/transaction")
}
