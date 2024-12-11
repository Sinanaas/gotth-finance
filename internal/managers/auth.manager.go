package managers

import (
	"fmt"
	"net/http"

	"github.com/Sinanaas/gotth-financial-tracker/internal/initializers"
	"github.com/Sinanaas/gotth-financial-tracker/internal/models"
	"github.com/Sinanaas/gotth-financial-tracker/internal/utils"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type AuthManager struct {
	DB     *gorm.DB
	config *initializers.Config
}

func NewAuthManager(db *gorm.DB, conf *initializers.Config) *AuthManager {
	return &AuthManager{DB: db, config: conf}
}

func (am *AuthManager) SignUp(ctx *gin.Context) (ret bool) {
	var payload models.SignUpInput
	if err := ctx.ShouldBind(&payload); err != nil {
		return false
	}

	if payload.Email == "" || payload.Password == "" {
		return false
	}

	if !utils.ValidateEmail(payload.Email) {
		return false
	}

	var user models.User
	result := am.DB.Where("email = ?", payload.Email).First(&user)
	if result.Error == nil {
		return false
	}

	if payload.Password != payload.ConfirmPassword {
		return false
	}

	hashedPassword, err := utils.HashPassword(payload.Password)
	if err != nil {
		return false
	}

	user = models.User{
		Username: payload.Username,
		Password: hashedPassword,
		Email:    payload.Email,
	}

	result = am.DB.Create(&user)
	if result.Error != nil {
		return false
	}

	return true
}

func (am *AuthManager) Login(ctx *gin.Context) (ret bool) {
	var payload models.SignInInput
	if err := ctx.ShouldBind(&payload); err != nil {
		return false
	}

	if payload.Email == "" || payload.Password == "" {
		return false
	}

	if !utils.ValidateEmail(payload.Email) {
		return false
	}

	var user models.User
	result := am.DB.Where("email = ?", payload.Email).First(&user)
	if result.Error != nil {
		return false
	}

	if err := utils.VerifyPassword(user.Password, payload.Password); err != nil {
		return false
	}

	accessToken, err := utils.GenerateToken(am.config.AccessTokenExpiresIn, user.ID.String(), am.config.AccessTokenPrivateKey)
	if err != nil {
		return false
	}

	refreshToken, err := utils.GenerateToken(am.config.RefreshTokenExpiresIn, user.ID.String(), am.config.RefreshTokenPrivateKey)
	if err != nil {
		return false
	}

	// set session
	session := sessions.Default(ctx)
	session.Set("user_id", user.ID.String())
	if err := session.Save(); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "Failed to save session"})
		return
	}

	ctx.SetCookie("access_token", accessToken, am.config.AccessTokenMaxAge*60, "/", "", false, true)
	ctx.SetCookie("refresh_token", refreshToken, am.config.RefreshTokenMaxAge*60, "/", "", false, true)
	ctx.SetCookie("logged_in", "true", am.config.AccessTokenMaxAge, "/", "", false, false)

	return true
}

func (am *AuthManager) RefreshToken(ctx *gin.Context) bool {
	cookie, err := ctx.Cookie("refresh_token")
	if err != nil {
		return false
	}
	config, _ := initializers.LoadConfig(".")
	sub, err := utils.ValidateToken(cookie, config.RefreshTokenPublicKey)
	if err != nil {
		return false
	}

	var user models.User
	result := am.DB.First(&user, "id = ?", fmt.Sprint(sub))
	if result.Error != nil {
		return false
	}
	accessToken, err := utils.GenerateToken(config.AccessTokenExpiresIn, user.ID.String(), config.AccessTokenPrivateKey)
	if err != nil {
		return false
	}

	ctx.SetCookie("access_token", accessToken, config.AccessTokenMaxAge, "/", "localhost", false, true)
	ctx.SetCookie("logged_in", "true", config.AccessTokenMaxAge, "/", "localhost", false, false)
	return true
}

func (am *AuthManager) Logout(ctx *gin.Context) {
	ctx.SetCookie("access_token", "", -1, "/", "", false, true)
	ctx.SetCookie("refresh_token", "", -1, "/", "", false, true)
	ctx.SetCookie("logged_in", "", -1, "/", "", false, false)

	session := sessions.Default(ctx)
	session.Clear()
}
