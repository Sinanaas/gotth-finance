package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/Sinanaas/gotth-financial-tracker/internal/controllers"
	"github.com/Sinanaas/gotth-financial-tracker/internal/utils"
	"github.com/gin-gonic/gin"
)

type DeleteCategoryHandler struct {
	BC *controllers.BasicController
}

func NewDeleteCategoryHandler(bc *controllers.BasicController) *DeleteCategoryHandler {
	return &DeleteCategoryHandler{BC: bc}
}

func (h *DeleteCategoryHandler) ServeHTTP(ctx *gin.Context) {
	userId := utils.GetSessionUserID(ctx)
	categoryId := ctx.Param("id")

	if err := h.BC.DeleteUserCategory(categoryId, userId); err != nil {
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
