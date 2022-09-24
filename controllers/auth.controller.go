package controllers

import (
	"github.com/akunsecured/emezen_api/models"
	"github.com/akunsecured/emezen_api/services"
	"github.com/akunsecured/emezen_api/utils"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator"
	"net/http"
)

type AuthController struct {
	authService services.AuthService
}

var validate = validator.New()

func NewAuthController(authService services.AuthService) AuthController {
	return AuthController{
		authService: authService,
	}
}

func (ac *AuthController) Register(ctx *gin.Context) {
	var credentials models.UserCredentials
	if err := ctx.ShouldBindJSON(&credentials); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	if err := validate.Struct(&credentials); err != nil {
		err = utils.ErrInvalidCredentialsFormat
		ctx.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	err := ac.authService.Register(&credentials)
	if err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"message": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "Successfully registered"})
}

func (ac *AuthController) Login(ctx *gin.Context) {
	var credentials models.UserCredentials
	if err := ctx.ShouldBindJSON(&credentials); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	if err := validate.Struct(&credentials); err != nil {
		err = utils.ErrInvalidCredentialsFormat
		ctx.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	token, err := ac.authService.Login(&credentials)
	if err != nil {
		switch err {
		case utils.ErrInvalidPassword:
			ctx.JSON(http.StatusUnauthorized, gin.H{"message": err.Error()})
		case utils.ErrNoAccountWithThisEmail:
			ctx.JSON(http.StatusNotFound, gin.H{"message": err.Error()})
		default:
			ctx.JSON(http.StatusBadGateway, gin.H{"message": err.Error()})
		}
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": token})
}

func (ac *AuthController) Update(ctx *gin.Context) {
	var credentials models.UserCredentials
	if err := ctx.ShouldBindJSON(&credentials); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	if err := validate.Struct(&credentials); err != nil {
		err = utils.ErrInvalidCredentialsFormat
		ctx.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	err := ac.authService.Update(&credentials)
	if err != nil {
		switch err {
		case utils.ErrNotExists:
			ctx.JSON(http.StatusUnprocessableEntity, gin.H{"message": err.Error()})
		default:
			ctx.JSON(http.StatusBadGateway, gin.H{"message": err.Error()})
		}
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "Successfully updated"})
}

func (ac *AuthController) RegisterAuthRoutes(rg *gin.RouterGroup) {
	authRoute := rg.Group("/auth")
	authRoute.POST("/register", ac.Register)
	authRoute.POST("/login", ac.Login)
	authRoute.PUT("/update", ac.Update)
}