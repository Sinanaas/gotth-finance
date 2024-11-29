package handlers

import (
	"github.com/Sinanaas/gotth-financial-tracker/internal/templates"
	"github.com/gin-gonic/gin"
	"net/http"
)

type GetLoginHandler struct{}

func NewGetLoginHandler() *GetLoginHandler {
	return &GetLoginHandler{}
}

func (h *GetLoginHandler) ServeHTTP(ctx *gin.Context) {
	c := templates.Login("Login")
	err := templates.Layout(c, "").Render(ctx.Request.Context(), ctx.Writer)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
}
