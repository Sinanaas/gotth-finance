package main

import (
	"log"

	"github.com/Sinanaas/gotth-financial-tracker/internal/controllers"
	"github.com/Sinanaas/gotth-financial-tracker/internal/initializers"
	"github.com/Sinanaas/gotth-financial-tracker/internal/managers"
	"github.com/Sinanaas/gotth-financial-tracker/internal/models"
	"github.com/Sinanaas/gotth-financial-tracker/internal/routers"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"github.com/go-co-op/gocron/v2"
	"github.com/Sinanaas/gotth-financial-tracker/internal/seeders"
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
		log.Fatal("❌ Could not load environment variables:", err)
	}

	initializers.ConnectDB(&config)
	
	initializers.DB.AutoMigrate(
		&models.User{},
  		&models.Category{},
        &models.Recurring{},
        &models.Transaction{},
        &models.Account{},
        &models.Loan{},
	)
	
	seeders.SeedCategories(initializers.DB)

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

	goCRON, err = gocron.NewScheduler()
	if err != nil {
		log.Fatal("❌ Could not create goCRON scheduler", err)
	}

	if err := BasicManager.LoadAndScheduleJobs(); err != nil {
		log.Fatal("❌ Could not load and schedule jobs", err)
	}

	goCRON.Start()
}

func main() {
	log.Fatal(server.Run(":" + config.ServerPort))
}
