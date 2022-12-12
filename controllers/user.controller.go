package controllers

import (
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/akunsecured/emezen_api/security"
	"github.com/akunsecured/emezen_api/utils"
	"github.com/form3tech-oss/jwt-go"

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

func (uc *UserController) CheckHeaderAuthorization(ctx *gin.Context) (*jwt.MapClaims, error) {
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

func (uc *UserController) GetUser(ctx *gin.Context) {
	_, err := uc.CheckHeaderAuthorization(ctx)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"message": err.Error()})
		return
	}

	userId := ctx.Param("id")
	user, err := uc.userService.GetUser(&userId)
	if err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"message": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": user})
}

func (uc *UserController) UpdateUser(ctx *gin.Context) {
	_, err := uc.CheckHeaderAuthorization(ctx)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"message": err.Error()})
		return
	}

	var user models.User
	if err := ctx.ShouldBindJSON(&user); err != nil {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"message": err.Error()})
		return
	}

	err = uc.userService.UpdateUser(&user)
	if err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"message": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "Success"})
}

func (uc *UserController) DeleteUser(ctx *gin.Context) {
	claims, err := uc.CheckHeaderAuthorization(ctx)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"message": err.Error()})
		return
	}

	userId := (*claims)["sub"].(string)
	err = uc.userService.DeleteUser(&userId)
	if err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"message": err.Error()})
		return
	}
	ctx.JSON(http.StatusNoContent, nil)
}

func (uc *UserController) UploadProfilePicture(ctx *gin.Context) {
	_, err := uc.CheckHeaderAuthorization(ctx)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"message": err.Error()})
		return
	}

	file, header, err := ctx.Request.FormFile("profile_picture")
	if err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"message": err.Error()})
		return
	}

	if file == nil || header == nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"message": "No files were uploaded."})
		return
	}

	fileName := header.Filename
	fmt.Println("File with name " + fileName + " is arrived.")

	filePath := "images/profile_pictures/" + fileName

	out, err := os.Create(filePath)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	defer out.Close()

	_, err = io.Copy(out, file)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": true})
}

func (uc *UserController) GetProfilePicture(ctx *gin.Context) {
	userId := ctx.Param("id")

	ctx.File("images/profile_pictures/" + userId + ".png")
}

func (uc *UserController) RegisterUserRoutes(rg *gin.RouterGroup) {
	userRoute := rg.Group("/user")
	userRoute.GET("/get/:id", uc.GetUser)
	userRoute.PUT("/update", uc.UpdateUser)
	userRoute.DELETE("/delete", uc.DeleteUser)
	userRoute.POST("/image/upload", uc.UploadProfilePicture)
	userRoute.GET("/image/:id", uc.GetProfilePicture)
}
