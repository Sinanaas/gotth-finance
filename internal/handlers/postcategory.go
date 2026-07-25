package handlers

import (
	"github.com/Sinanaas/gotth-financial-tracker/internal/managers"
	"github.com/Sinanaas/gotth-financial-tracker/internal/templates"
	"github.com/Sinanaas/gotth-financial-tracker/internal/utils"
	"github.com/Sinanaas/gotth-financial-tracker/internal/validation"
	"github.com/gin-gonic/gin"
)

type PostCategoryHandler struct {
	BM *managers.BasicManager
}

func NewPostCategoryHandler(bm *managers.BasicManager) *PostCategoryHandler {
	return &PostCategoryHandler{BM: bm}
}

func (h *PostCategoryHandler) ServeHTTP(ctx *gin.Context) {
	userId := utils.GetSessionUserID(ctx)

	name, err := validation.Required(ctx.PostForm("Name"), "category name")
	if err != nil {
		swalError(ctx, err.Error())
		return
	}

	if err := h.BM.CreateUserCategory(userId, name, ctx.PostForm("Description")); err != nil {
		swalError(ctx, "Failed to create category")
		return
	}

	categories, err := h.BM.GetUserCategories(userId)
	if err != nil {
		swalError(ctx, "Failed to reload categories")
		return
	}
	ctx.Status(200)
	_ = templates.CategoriesPanel(categories).Render(ctx.Request.Context(), ctx.Writer)
}
