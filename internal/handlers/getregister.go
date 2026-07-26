package handlers

import (
	"net/http"

	"github.com/Sinanaas/gotth-financial-tracker/internal/controllers"
	"github.com/Sinanaas/gotth-financial-tracker/internal/templates"
	"github.com/gin-gonic/gin"
)

type GetRegisterHandler struct {
	AC *controllers.AuthController
}

func NewGetRegisterHandler(ac *controllers.AuthController) *GetRegisterHandler {
	return &GetRegisterHandler{AC: ac}
}

func (h *GetRegisterHandler) ServeHTTP(ctx *gin.Context) {
	// When registration is closed, don't serve the signup form.
	if !h.AC.RegistrationEnabled() {
		ctx.Redirect(http.StatusFound, "/login")
		return
	}

	c := templates.Register("Register")
	if err := templates.AuthLayout(c).Render(ctx.Request.Context(), ctx.Writer); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}
}
