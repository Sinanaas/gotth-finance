package handlers

import (
	"github.com/Sinanaas/gotth-financial-tracker/internal/controllers"
	"github.com/Sinanaas/gotth-financial-tracker/internal/templates"
	"github.com/gin-gonic/gin"
)

type GetLogoutHandler struct {
	AC *controllers.AuthController
}

func NewGetLogoutHandler(ac *controllers.AuthController) *GetLogoutHandler {
	return &GetLogoutHandler{AC: ac}
}

func (h *GetLogoutHandler) ServeHTTP(c *gin.Context) {
	confirm := c.Query("confirm")
	if confirm == "true" {
		h.AC.Logout(c)
		return
	}

	err := templates.Modal("Logout", "Are you sure you want to logout?").Render(c.Request.Context(), c.Writer)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	h.AC.Logout(c)
}
