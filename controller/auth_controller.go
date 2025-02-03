package controller

import (
	"project_rentalmobil/model"
	"project_rentalmobil/usecase"

	"github.com/gin-gonic/gin"
)

type AuthController struct {
	authUc usecase.AuthenticationUseCase
	rg     *gin.RouterGroup
}

func (a *AuthController) Route() {
	a.rg.POST("/login", a.loginHandler)
}

func (a *AuthController) loginHandler(ctx *gin.Context) {
	var payload model.UserCredential

	err := ctx.ShouldBindJSON(&payload)
	if err != nil {
		ctx.JSON(400, gin.H{"err": err})
		return
	}

	token, err := a.authUc.Login(payload.Username, payload.Password)
	if err != nil {
		ctx.JSON(500, gin.H{"err": err})
		return
	}

	ctx.JSON(201, gin.H{"token": token})
}

func NewAuthController(authUc usecase.AuthenticationUseCase, rg *gin.RouterGroup) *AuthController {
	return &AuthController{authUc: authUc, rg: rg}
}
