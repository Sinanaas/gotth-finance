package routers

import (
	"github.com/Sinanaas/gotth-financial-tracker/internal/controllers"
	"github.com/Sinanaas/gotth-financial-tracker/internal/handlers"
	"github.com/Sinanaas/gotth-financial-tracker/internal/middleware"
	"github.com/gin-gonic/gin"
)

type BasicRouter struct {
	BC *controllers.BasicController
	RG *gin.RouterGroup
}

func NewBasicRouter(bc *controllers.BasicController, rg *gin.RouterGroup) *BasicRouter {
	return &BasicRouter{BC: bc, RG: rg}
}

func (br *BasicRouter) BasicRoute(rg *gin.RouterGroup, bc *controllers.BasicController) {
	rg.GET("/", middleware.DeserializeUser(), handlers.NewGetHomeHandler(bc).ServeHTTP)
	rg.GET("/transaction", middleware.DeserializeUser(), handlers.NewGetTransaction(bc).ServeHTTP)
	rg.POST("/transaction", middleware.DeserializeUser(), handlers.NewPostTransactionHandler(bc).ServeHTTP)
	rg.GET("/accounts", middleware.DeserializeUser(), handlers.NewGetAccountsHandler(bc).ServeHTTP)
	rg.POST("/account", middleware.DeserializeUser(), handlers.NewPostAccountHandler(bc).ServeHTTP)
	rg.GET("/account/balance", middleware.DeserializeUser(), handlers.NewGetBalanceHandler(bc).ServeHTTP)
	rg.GET("/recurring", middleware.DeserializeUser(), handlers.NewGetRecurringHandler(bc).ServeHTTP)
	rg.POST("/recurring", middleware.DeserializeUser(), handlers.NewPostRecurringHandler(bc).ServeHTTP)
	rg.GET("/loans", middleware.DeserializeUser(), handlers.NewGetLoanHandler(bc).ServeHTTP)
	rg.POST("/loan", middleware.DeserializeUser(), handlers.NewPostLoanHandler(bc).ServeHTTP)
	rg.POST("/loan/finish", middleware.DeserializeUser(), handlers.NewPostFinishLoanHandler(bc).ServeHTTP)
}
