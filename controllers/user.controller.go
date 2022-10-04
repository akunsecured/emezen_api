package controllers

import (
	"github.com/akunsecured/emezen_api/security"
	"github.com/akunsecured/emezen_api/utils"
	"net/http"

	"github.com/akunsecured/emezen_api/models"
	"github.com/akunsecured/emezen_api/services"
	"github.com/gin-gonic/gin"
)

type UserController struct {
	userService services.UserService
}

func NewUserController(userService services.UserService) UserController {
	return UserController{
		userService: userService,
	}
}

func (uc *UserController) GetUser(ctx *gin.Context) {
	tokenStr := ctx.GetHeader("Authorization")
	if tokenStr == "" {
		ctx.JSON(http.StatusUnauthorized, gin.H{"message": utils.ErrMissingAuthToken})
		return
	}

	claims, err := security.ParseToken(tokenStr)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"message": err.Error()})
		return
	}

	userId := (*claims)["sub"].(string)
	user, err := uc.userService.GetUser(&userId)
	if err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"message": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, user)
}

func (uc *UserController) UpdateUser(ctx *gin.Context) {
	var user models.User
	if err := ctx.ShouldBindJSON(&user); err != nil {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"message": err.Error()})
		return
	}
	err := uc.userService.UpdateUser(&user)
	if err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"message": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "Success"})
}

func (uc *UserController) DeleteUser(ctx *gin.Context) {
	userName := ctx.Param("name")
	err := uc.userService.DeleteUser(&userName)
	if err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"message": err.Error()})
		return
	}
	ctx.JSON(http.StatusNoContent, nil)
}

func (uc *UserController) RegisterUserRoutes(rg *gin.RouterGroup) {
	userRoute := rg.Group("/user")
	userRoute.GET("/get", uc.GetUser)
	userRoute.PUT("/update", uc.UpdateUser)
	userRoute.DELETE("/delete/:name", uc.DeleteUser)
}
