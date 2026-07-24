package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/Sinanaas/gotth-financial-tracker/internal/controllers"
	"github.com/Sinanaas/gotth-financial-tracker/internal/utils"
	"github.com/gin-gonic/gin"
)

type PostCategoryHandler struct {
	BC *controllers.BasicController
}

func NewPostCategoryHandler(bc *controllers.BasicController) *PostCategoryHandler {
	return &PostCategoryHandler{BC: bc}
}

func (h *PostCategoryHandler) ServeHTTP(ctx *gin.Context) {
	userId := utils.GetSessionUserID(ctx)

	name := ctx.PostForm("Name")
	description := ctx.PostForm("Description")

	if name == "" {
		ctx.Status(http.StatusBadRequest)
		return
	}

	if err := h.BC.CreateUserCategory(userId, name, description); err != nil {
		swalData, _ := json.Marshal(map[string]interface{}{
			"swal:alert": map[string]interface{}{
				"title":    "Error",
				"text":     "Failed to create category",
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
