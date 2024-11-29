package routers

import (
	"github.com/Sinanaas/gotth-financial-tracker/internal/controllers"
	"github.com/Sinanaas/gotth-financial-tracker/internal/handlers"
	"github.com/Sinanaas/gotth-financial-tracker/internal/middleware"
	"github.com/gin-gonic/gin"
)

type AuthRouter struct {
	authController controllers.AuthController
}

func NewAuthRouter(authController controllers.AuthController) AuthRouter {
	return AuthRouter{authController}
}

func (ar *AuthRouter) AuthRoute(rg *gin.RouterGroup) {
	rg.GET("/login", handlers.NewGetLoginHandler().ServeHTTP)
	rg.POST("/login", ar.authController.SignIn)
	rg.GET("/register", handlers.NewGetRegisterHandler().ServeHTTP)
	rg.POST("/register", ar.authController.SignUp)
	rg.GET("/logout", middleware.DeserializeUser(), ar.authController.Logout)
}
