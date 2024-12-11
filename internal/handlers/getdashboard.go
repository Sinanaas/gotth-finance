package handlers

import (
	"github.com/Sinanaas/gotth-financial-tracker/internal/controllers"
	"github.com/gin-gonic/gin"
)

type GetDashboardHandler struct {
	BC *controllers.BasicController
}

func NewGetDashboardHandler(bc *controllers.BasicController) *GetDashboardHandler {
	return &GetDashboardHandler{BC: bc}
}

func (h *GetDashboardHandler) ServeHTTP(ctx *gin.Context) {

}
