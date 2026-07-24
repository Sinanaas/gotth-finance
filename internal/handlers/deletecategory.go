package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/Sinanaas/gotth-financial-tracker/internal/managers"
	"github.com/Sinanaas/gotth-financial-tracker/internal/utils"
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
	categoryId := ctx.Param("id")

	if err := h.BM.DeleteUserCategory(categoryId, userId); err != nil {
		swalData, _ := json.Marshal(map[string]interface{}{
			"swal:alert": map[string]interface{}{
				"title":    "Error",
				"text":     "Failed to delete category",
				"icon":     "error",
				"redirect": "/categories",
			},
		})
		ctx.Header("HX-Trigger", string(swalData))
		ctx.Status(http.StatusInternalServerError)
		return
	}

	ctx.Header("HX-Redirect", "/categories")
	ctx.Status(http.StatusOK)
}
