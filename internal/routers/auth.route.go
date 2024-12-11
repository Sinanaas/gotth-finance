package routers

import (
	"github.com/Sinanaas/gotth-financial-tracker/internal/controllers"
	"github.com/Sinanaas/gotth-financial-tracker/internal/handlers"
	"github.com/Sinanaas/gotth-financial-tracker/internal/middleware"
	"github.com/gin-gonic/gin"
)

type AuthRouter struct {
	AC *controllers.AuthController
	RG *gin.RouterGroup
}

func NewAuthRouter(ac *controllers.AuthController, rg *gin.RouterGroup) *AuthRouter {
	return &AuthRouter{AC: ac, RG: rg}
}

func (ar *AuthRouter) AuthRoute(rg *gin.RouterGroup, ac *controllers.AuthController) {
	rg.GET("/login", handlers.NewGetLoginHandler().ServeHTTP)
	rg.POST("/login", handlers.NewPostLoginHandler(ac).ServeHTTP)
	rg.GET("/register", handlers.NewGetRegisterHandler().ServeHTTP)
	rg.POST("/register", handlers.NewPostRegisterHandler(ac).ServeHTTP)
	rg.GET("/logout", middleware.DeserializeUser(), handlers.NewGetLogoutHandler(ac).ServeHTTP)
}
