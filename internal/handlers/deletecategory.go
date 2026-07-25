package handlers

import (
	"github.com/Sinanaas/gotth-financial-tracker/internal/managers"
	"github.com/Sinanaas/gotth-financial-tracker/internal/templates"
	"github.com/Sinanaas/gotth-financial-tracker/internal/utils"
	"github.com/Sinanaas/gotth-financial-tracker/internal/validation"
	"github.com/gin-gonic/gin"
)

type DeleteCategoryHandler struct {
	BM *managers.BasicManager
}

func NewDeleteCategoryHandler(bm *managers.BasicManager) *DeleteCategoryHandler {
	return &DeleteCategoryHandler{BM: bm}
}

func (h *DeleteCategoryHandler) ServeHTTP(ctx *gin.Context) {
	userId := utils.GetSessionUserID(ctx)

	categoryId, err := validation.UUIDField(ctx.Param("id"), "category")
	if err != nil {
		swalError(ctx, err.Error())
		return
	}

	if err := h.BM.DeleteUserCategory(categoryId, userId); err != nil {
		swalError(ctx, "Failed to delete category")
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
