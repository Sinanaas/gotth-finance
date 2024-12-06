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
	rg.GET("/accounts", middleware.DeserializeUser(), handlers.NewGetAccountsHandler)
	rg.POST("/account", middleware.DeserializeUser(), br.basicController.CreateAccount)
	rg.GET("/account/balance", middleware.DeserializeUser(), br.basicController.GetAccountBalance)
	rg.GET("/recurring", middleware.DeserializeUser(), handlers.NewGetRecurringHandler().ServeHTTP)
	rg.POST("/recurring", middleware.DeserializeUser(), br.basicController.CreateRecurring)

}
