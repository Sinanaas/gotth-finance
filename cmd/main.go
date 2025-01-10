package main

import (
	"github.com/Sinanaas/gotth-financial-tracker/internal/controllers"
	"github.com/Sinanaas/gotth-financial-tracker/internal/initializers"
	"github.com/Sinanaas/gotth-financial-tracker/internal/managers"
	"github.com/Sinanaas/gotth-financial-tracker/internal/routers"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"github.com/go-co-op/gocron/v2"
	"log"
)

var (
	server *gin.Engine
	router *gin.RouterGroup
	config initializers.Config
	err    error

	BasicController *controllers.BasicController
	BasicManager    *managers.BasicManager
	BasicRouter     *routers.BasicRouter

	AuthController *controllers.AuthController
	AuthManager    *managers.AuthManager
	AuthRouter     *routers.AuthRouter

	goCRON gocron.Scheduler
)

func init() {
	config, err = initializers.LoadConfig(".")
	if err != nil {
		log.Fatal("❌ Could not load environment variables", err)
	}

	goCRON, err = gocron.NewScheduler()
	if err != nil {
		log.Fatal("❌ Could not create goCRON scheduler", err)
	}

	initializers.ConnectDB(&config)

	AuthManager = managers.NewAuthManager(initializers.DB, &config)
	AuthController = controllers.NewAuthController(AuthManager)
	AuthRouter = routers.NewAuthRouter(AuthController, router)

	BasicManager = managers.NewBasicManager(initializers.DB, &goCRON)
	BasicController = controllers.NewBasicController(BasicManager)
	BasicRouter = routers.NewBasicRouter(BasicController, router)

	server = gin.Default()
	store := cookie.NewStore([]byte(config.SessionSecretKey))
	server.Use(sessions.Sessions("mysession", store))

	router = server.Group("/")

	BasicRouter.BasicRoute(router, BasicController)
	AuthRouter.AuthRoute(router, AuthController)

	goCRON.Start()
}

func main() {
	// print all active gocron jobs
	for _, job := range goCRON.Jobs() {
		log.Println("JOB:", job)
	}
	log.Println("HELLO WORLD!")

	log.Fatal(server.Run(":" + config.ServerPort))
}
