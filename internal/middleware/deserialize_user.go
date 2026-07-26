package middleware

import (
	"fmt"
	"github.com/Sinanaas/gotth-financial-tracker/internal/controllers"
	"github.com/Sinanaas/gotth-financial-tracker/internal/initializers"
	"github.com/Sinanaas/gotth-financial-tracker/internal/managers"
	"github.com/Sinanaas/gotth-financial-tracker/internal/models"
	"github.com/Sinanaas/gotth-financial-tracker/internal/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

func DeserializeUser() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		cookie, err := ctx.Cookie("access_token")
		if err != nil || cookie == "" {
			ctx.Redirect(http.StatusSeeOther, "/login")
			ctx.Abort()
			return
		}

		config, _ := initializers.LoadConfig(".")
		sub, err := utils.ValidateToken(cookie, config.AccessTokenPublicKey)
		if err != nil {
			ctx.Redirect(http.StatusSeeOther, "/login")
			ctx.Abort()
			return
		}

		var user models.User
		if err := initializers.DB.First(&user, "id = ?", fmt.Sprint(sub)).Error; err != nil {
			ctx.Redirect(http.StatusSeeOther, "/login")
			ctx.Abort()
			return
		}

		am := managers.NewAuthManager(initializers.DB, &config)
		at := controllers.NewAuthController(am)
		at.RefreshToken(ctx)

		ctx.Set("currentUser", user)
		ctx.Next()
	}
}
