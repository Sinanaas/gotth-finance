package handlers

import (
	"net/http"
	"strconv"

	"github.com/Sinanaas/gotth-financial-tracker/internal/managers"
	"github.com/Sinanaas/gotth-financial-tracker/internal/models"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

type PostTransferHandler struct {
	BM *managers.BasicManager
}

func NewPostTransferHandler(bm *managers.BasicManager) *PostTransferHandler {
	return &PostTransferHandler{BM: bm}
}

func (h *PostTransferHandler) ServeHTTP(ctx *gin.Context) {
	var payload models.TransferRequest
	var err error

	payload.FromAccountID = ctx.PostForm("FromAccount")
	payload.ToAccountID = ctx.PostForm("ToAccount")
	payload.Amount, err = strconv.ParseFloat(ctx.PostForm("Amount"), 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	payload.Date = ctx.PostForm("Date")
	payload.Description = ctx.PostForm("Description")

	session := sessions.Default(ctx)
	var userId string
	v := session.Get("user_id")
	if v != nil {
		userId = v.(string)
	}
	payload.UserID = userId

	if err := h.BM.CreateTransfer(payload); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.Header("HX-Redirect", "/accounts")
}
