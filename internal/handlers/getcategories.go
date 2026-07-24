package handlers

import (
	"net/http"

	"github.com/Sinanaas/gotth-financial-tracker/internal/controllers"
	"github.com/Sinanaas/gotth-financial-tracker/internal/templates"
	"github.com/Sinanaas/gotth-financial-tracker/internal/utils"
	"github.com/gin-gonic/gin"
)

type GetCategoriesHandler struct {
	BC *controllers.BasicController
}

func NewGetCategoriesHandler(bc *controllers.BasicController) *GetCategoriesHandler {
	return &GetCategoriesHandler{BC: bc}
}

func (h *GetCategoriesHandler) ServeHTTP(ctx *gin.Context) {
	cookie, _ := ctx.Cookie("access_token")
	userId := utils.GetSessionUserID(ctx)

	categories, err := h.BC.GetUserCategories(userId)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to load categories"})
		return
	}

	c := templates.Categories(categories)
	if err := templates.Layout(c, cookie).Render(ctx.Request.Context(), ctx.Writer); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to render template"})
	}
}
