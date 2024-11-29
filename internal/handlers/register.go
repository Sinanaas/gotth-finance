package handlers

import (
	"github.com/Sinanaas/gotth-financial-tracker/internal/templates"
	"github.com/gin-gonic/gin"
	"net/http"
)

type GetRegisterHandler struct{}

func NewGetRegisterHandler() *GetRegisterHandler {
	return &GetRegisterHandler{}
}

func (h *GetRegisterHandler) ServeHTTP(ctx *gin.Context) {
	c := templates.Register("Register")
	err := templates.Layout(c, "").Render(ctx.Request.Context(), ctx.Writer)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
}
