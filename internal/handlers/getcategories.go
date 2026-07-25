package handlers

import (
	"github.com/Sinanaas/gotth-financial-tracker/internal/managers"
	"github.com/Sinanaas/gotth-financial-tracker/internal/templates"
	"github.com/Sinanaas/gotth-financial-tracker/internal/utils"
	"github.com/gin-gonic/gin"
)

type GetCategoriesHandler struct {
	BM *managers.BasicManager
}

func NewGetCategoriesHandler(bm *managers.BasicManager) *GetCategoriesHandler {
	return &GetCategoriesHandler{BM: bm}
}

func (h *GetCategoriesHandler) ServeHTTP(ctx *gin.Context) {
	cookie, _ := ctx.Cookie("access_token")
	userId := utils.GetSessionUserID(ctx)

	categories, err := h.BM.GetUserCategories(userId)
	if err != nil {
		renderErrorPage(ctx, cookie, "Failed to load categories")
		return
	}

	c := templates.Categories(categories)
	if err := templates.Layout(c, cookie).Render(ctx.Request.Context(), ctx.Writer); err != nil {
		renderErrorPage(ctx, cookie, "Failed to render template")
	}
}
