package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Sinanaas/gotth-financial-tracker/internal/controllers"
	"github.com/Sinanaas/gotth-financial-tracker/internal/initializers"
	"github.com/Sinanaas/gotth-financial-tracker/internal/managers"
	"github.com/Sinanaas/gotth-financial-tracker/internal/models"
	"github.com/Sinanaas/gotth-financial-tracker/internal/routers"
	"github.com/Sinanaas/gotth-financial-tracker/internal/seeders"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"github.com/go-co-op/gocron/v2"
)

var (
	server *gin.Engine
	router *gin.RouterGroup
	config initializers.Config
	err    error

	BasicManager *managers.BasicManager
	BasicRouter  *routers.BasicRouter

	AuthController *controllers.AuthController
	AuthManager    *managers.AuthManager
	AuthRouter     *routers.AuthRouter

	goCRON gocron.Scheduler
)

func init() {
	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, nil)))

	config, err = initializers.LoadConfig(".")
	if err != nil {
		slog.Error("could not load environment variables", "err", err)
		os.Exit(1)
	}
	if err := config.Validate(); err != nil {
		slog.Error("invalid configuration", "err", err)
		os.Exit(1)
	}

	initializers.ConnectDB(&config)

	initializers.DB.AutoMigrate(
		&models.User{},
		&models.Category{},
		&models.Recurring{},
		&models.Transaction{},
		&models.Account{},
		&models.Loan{},
		&models.Budget{},
	)

	seeders.SeedCategories(initializers.DB)

	AuthManager = managers.NewAuthManager(initializers.DB, &config)
	AuthController = controllers.NewAuthController(AuthManager)
	AuthRouter = routers.NewAuthRouter(AuthController, router)

	BasicManager = managers.NewBasicManager(initializers.DB, &goCRON)
	BasicRouter = routers.NewBasicRouter(BasicManager, router)

	server = gin.Default()
	store := cookie.NewStore([]byte(config.SessionSecretKey))
	store.Options(sessions.Options{Path: "/", MaxAge: 86400 * 7, HttpOnly: true, Secure: config.CookieSecure, SameSite: http.SameSiteLaxMode})
	server.Use(sessions.Sessions("mysession", store))

	router = server.Group("/")

	BasicRouter.BasicRoute(router, BasicManager)
	AuthRouter.AuthRoute(router, AuthController)

	goCRON, err = gocron.NewScheduler()
	if err != nil {
		slog.Error("could not create scheduler", "err", err)
		os.Exit(1)
	}

	if err := BasicManager.LoadAndScheduleJobs(); err != nil {
		slog.Error("could not load and schedule jobs", "err", err)
		os.Exit(1)
	}

	goCRON.Start()
}

func main() {
	srv := &http.Server{Addr: ":" + config.ServerPort, Handler: server}

	go func() {
		slog.Info("server started", "port", config.ServerPort)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			slog.Error("server error", "err", err)
			os.Exit(1)
		}
	}()

	// Block until an interrupt/terminate signal, then drain gracefully.
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	<-ctx.Done()
	stop()
	slog.Info("shutting down")

	_ = goCRON.Shutdown()

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := srv.Shutdown(shutdownCtx); err != nil {
		slog.Error("graceful shutdown failed", "err", err)
	}
}
