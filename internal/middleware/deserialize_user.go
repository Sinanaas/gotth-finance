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
		var accessToken string
		cookie, err := ctx.Cookie("access_token")
		if err != nil {
			ctx.Redirect(http.StatusSeeOther, "/login")
			return
		}

		accessToken = cookie
		if accessToken == "" {
			ctx.Redirect(http.StatusSeeOther, "/login")
			return
		}

		config, _ := initializers.LoadConfig(".")
		sub, err := utils.ValidateToken(accessToken, config.AccessTokenPublicKey)

		if err != nil {
			ctx.Redirect(http.StatusSeeOther, "/login")
			return
		}

		var user models.User
		result := initializers.DB.First(&user, "id = ?", fmt.Sprint(sub))
		if result.Error != nil {
			ctx.Redirect(http.StatusSeeOther, "/login")
			return
		}

		am := managers.NewAuthManager(initializers.DB, &config)
		at := controllers.NewAuthController(am)
		at.RefreshToken(ctx)

		ctx.Set("currentUser", user)
		ctx.Next()
	}
}
