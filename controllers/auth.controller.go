package controllers

import (
	"github.com/akunsecured/emezen_api/models"
	"github.com/akunsecured/emezen_api/security"
	"github.com/akunsecured/emezen_api/services"
	"github.com/akunsecured/emezen_api/utils"
	"github.com/form3tech-oss/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator"
	"log"
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

func (ac *AuthController) CheckHeaderAuthorization(ctx *gin.Context) (*jwt.MapClaims, error) {
	tokenStr := ctx.GetHeader("Authorization")
	if tokenStr == "" {
		return nil, utils.ErrMissingAuthToken
	}

	claims, err := security.ParseToken(tokenStr)
	if err != nil {
		return nil, err
	}

	return claims, nil
}

func (ac *AuthController) Register(ctx *gin.Context) {
	var userDataWithCredentials models.UserDataWithCredentials
	if err := ctx.ShouldBindJSON(&userDataWithCredentials); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	if err := validate.Struct(&userDataWithCredentials); err != nil {
		log.Print(err)
		err = utils.ErrInvalidCredentialsFormat
		ctx.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	wrappedToken, err := ac.authService.Register(&userDataWithCredentials)
	if err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"message": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": wrappedToken})
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
	wrappedToken, err := ac.authService.Login(&credentials)
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
	ctx.JSON(http.StatusOK, gin.H{"message": wrappedToken})
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

func (ac *AuthController) RefreshToken(ctx *gin.Context) {
	claims, err := ac.CheckHeaderAuthorization(ctx)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"message": err.Error()})
		return
	}

	newToken, err := ac.authService.NewAccessToken(claims)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": newToken})
}

func (ac *AuthController) CurrentUser(ctx *gin.Context) {
	claims, err := ac.CheckHeaderAuthorization(ctx)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"message": err.Error()})
		return
	}

	currentUser, err := ac.authService.CurrentUser(claims)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": currentUser})
}

func (ac *AuthController) RegisterAuthRoutes(rg *gin.RouterGroup) {
	authRoute := rg.Group("/auth")
	authRoute.POST("/register", ac.Register)
	authRoute.POST("/login", ac.Login)
	authRoute.PUT("/update", ac.Update)
	authRoute.GET("/refresh", ac.RefreshToken)
	authRoute.GET("/current", ac.CurrentUser)
}
