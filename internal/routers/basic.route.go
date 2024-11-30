package routers

import (
	"github.com/Sinanaas/gotth-financial-tracker/internal/controllers"
	"github.com/Sinanaas/gotth-financial-tracker/internal/handlers"
	"github.com/Sinanaas/gotth-financial-tracker/internal/middleware"
	"github.com/gin-gonic/gin"
)

type BasicRouter struct {
	basicController controllers.BasicController
}

func NewBasicRouter(basicController controllers.BasicController) BasicRouter {
	return BasicRouter{basicController}
}

func (br *BasicRouter) BasicRoute(rg *gin.RouterGroup) {
	rg.GET("/", middleware.DeserializeUser(), handlers.NewGetHomeHandler().ServeHTTP)
	rg.GET("/transaction", middleware.DeserializeUser(), handlers.NewGetTransaction().ServeHTTP)
	rg.POST("/transaction", middleware.DeserializeUser(), br.basicController.CreateTransaction)
}
