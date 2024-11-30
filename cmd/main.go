package main

import (
	"github.com/Sinanaas/gotth-financial-tracker/internal/controllers"
	"github.com/Sinanaas/gotth-financial-tracker/internal/initializers"
	"github.com/Sinanaas/gotth-financial-tracker/internal/managers"
	"github.com/Sinanaas/gotth-financial-tracker/internal/routers"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"log"
)

var (
	server          *gin.Engine
	BasicController controllers.BasicController
	BasicManager    managers.BasicManager
	BasicRouter     routers.BasicRouter

	AuthController controllers.AuthController
	AuthManager    managers.AuthManager
	AuthRouter     routers.AuthRouter
)

func init() {
	config, err := initializers.LoadConfig(".")
	if err != nil {
		log.Fatal("❌ Could not load environment variables", err)
	}

	initializers.ConnectDB(&config)

	AuthManager = managers.NewAuthManager(initializers.DB, &config)
	AuthController = controllers.NewAuthController(AuthManager)
	AuthRouter = routers.NewAuthRouter(AuthController)

	BasicManager = managers.NewBasicManager(initializers.DB)
	BasicController = controllers.NewBasicController(BasicManager)
	BasicRouter = routers.NewBasicRouter(BasicController)

	server = gin.Default()
	store := cookie.NewStore([]byte(config.SessionSecretKey))
	server.Use(sessions.Sessions("mysession", store))
}

func main() {
	config, err := initializers.LoadConfig(".")
	if err != nil {
		log.Fatal("❌ Could not load environment variables", err)
	}

	router := server.Group("/")

	BasicRouter.BasicRoute(router)
	AuthRouter.AuthRoute(router)

	log.Fatal(server.Run(":" + config.ServerPort))
}
