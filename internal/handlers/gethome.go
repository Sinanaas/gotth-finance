package handlers

import (
	"github.com/Sinanaas/gotth-financial-tracker/internal/templates"
	"github.com/gin-gonic/gin"
	"net/http"
)

type GetHomeHandler struct {
}

func NewGetHomeHandler() *GetHomeHandler {
	return &GetHomeHandler{}
}

func (h *GetHomeHandler) ServeHTTP(ctx *gin.Context) {
	c := templates.Home()
	cookie, _ := ctx.Cookie("access_token")
	err := templates.Layout(c, cookie).Render(ctx.Request.Context(), ctx.Writer)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
}
