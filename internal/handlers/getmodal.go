package handlers

import (
	"github.com/gin-gonic/gin"
)

type GetModal struct{}

func NewGetModal() *GetModal {
	return &GetModal{}
}

func (h *GetModal) ServeHTTP(ctx *gin.Context) {
	//c := templates.Modal()
	//err := templates.Layout(c, "").Render(ctx.Request.Context(), ctx.Writer)
	//
	//if err != nil {
	//	ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	//	return
	//}
}
