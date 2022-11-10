package controllers

import (
	"net/http"

	"github.com/akunsecured/emezen_api/models"
	"github.com/akunsecured/emezen_api/security"
	"github.com/akunsecured/emezen_api/services"
	"github.com/akunsecured/emezen_api/utils"
	"github.com/form3tech-oss/jwt-go"
	"github.com/gin-gonic/gin"
)

type ProductController struct {
	productService services.ProductService
}

func NewProductController(productService services.ProductService) ProductController {
	return ProductController{
		productService: productService,
	}
}

func (pc *ProductController) CheckHeaderAuthorization(ctx *gin.Context) (*jwt.MapClaims, error) {
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

func (pc *ProductController) CreateProduct(ctx *gin.Context) {
	var product models.Product
	if err := ctx.ShouldBindJSON(&product); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	if err := validate.Struct(&product); err != nil {
		err = utils.ErrInvalidProductFormat
		ctx.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	productId, err := pc.productService.AddProduct(&product)
	if err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"message": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": productId})
}

func (pc *ProductController) GetProduct(ctx *gin.Context) {
	_, err := pc.CheckHeaderAuthorization(ctx)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"message": err.Error()})
		return
	}

	productId := ctx.Param("id")
	product, err := pc.productService.GetProduct(&productId)
	if err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"message": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": product})
}

func (pc *ProductController) GetAllProducts(ctx *gin.Context) {
	_, err := pc.CheckHeaderAuthorization(ctx)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"message": err.Error()})
		return
	}
	products, err := pc.productService.GetAllProducts()
	if err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"message": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": products})
}

func (pc *ProductController) UpdateProduct(ctx *gin.Context) {
	claims, err := pc.CheckHeaderAuthorization(ctx)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"message": err.Error()})
		return
	}

	productId := ctx.Param("id")
	product, err := pc.productService.GetProduct(&productId)
	if err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"message": err.Error()})
		return
	}
	userId := (*claims)["sub"].(string)

	if product.SellerID != userId {
		ctx.JSON(http.StatusUnauthorized, gin.H{"message": "tried to update other people's product"})
		return
	}

	// TODO: implement update
	ctx.JSON(http.StatusOK, gin.H{"message": "updated"})
}

func (pc *ProductController) DeleteProduct(ctx *gin.Context) {
	claims, err := pc.CheckHeaderAuthorization(ctx)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"message": err.Error()})
		return
	}

	productId := ctx.Param("id")
	product, err := pc.productService.GetProduct(&productId)
	if err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"message": err.Error()})
		return
	}
	userId := (*claims)["sub"].(string)

	if product.SellerID != userId {
		ctx.JSON(http.StatusUnauthorized, gin.H{"message": "tried to delete other people's product"})
		return
	}

	err = pc.productService.DeleteProduct(&productId)
	if err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"message": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "product with id " + productId + " is deleted"})
}

func (pc *ProductController) RegisterProductRoutes(rg *gin.RouterGroup) {
	productRoute := rg.Group("/product")
	productRoute.POST("/create", pc.CreateProduct)
	productRoute.GET("/get/:id", pc.GetProduct)
	productRoute.GET("/get_all", pc.GetAllProducts)
	productRoute.PUT("/update/:id", pc.UpdateProduct)
	productRoute.DELETE("/delete/:id", pc.DeleteProduct)
}
