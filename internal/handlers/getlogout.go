package handlers

import (
	"github.com/Sinanaas/gotth-financial-tracker/internal/controllers"
	"github.com/gin-gonic/gin"
)

type GetLogoutHandler struct {
	AC *controllers.AuthController
}

func NewGetLogoutHandler(ac *controllers.AuthController) *GetLogoutHandler {
	return &GetLogoutHandler{AC: ac}
}

func (h *GetLogoutHandler) ServeHTTP(c *gin.Context) {
	h.AC.Logout(c)
}
