package controllers

import (
	"github.com/Sinanaas/gotth-financial-tracker/internal/managers"
	"github.com/gin-gonic/gin"
)

type AuthController struct {
	AuthManager *managers.AuthManager
}

func NewAuthController(authManager *managers.AuthManager) *AuthController {
	return &AuthController{AuthManager: authManager}
}

func (ac *AuthController) SignUp(ctx *gin.Context) {
	ret := ac.AuthManager.SignUp(ctx)
	if ret {
		ctx.Header("HX-Redirect", "/login")
		ctx.Status(200)
	} else {
		ctx.Header("HX-Redirect", "/register")
		ctx.Status(401)
	}
}

func (ac *AuthController) SignIn(ctx *gin.Context) {
	ret := ac.AuthManager.Login(ctx)
	if ret {
		ctx.Header("HX-Redirect", "/")
		ctx.Status(200)
	} else {
		ctx.Header("HX-Redirect", "/login")
		ctx.Status(401)
	}
}

func (ac *AuthController) RefreshToken(ctx *gin.Context) {
	if ac.AuthManager.RefreshToken(ctx) == false {
		ctx.Header("HX-Redirect", "/login")
	}
}

func (ac *AuthController) Logout(ctx *gin.Context) {
	ac.AuthManager.Logout(ctx)
	ctx.Header("HX-Redirect", "/login")
}
