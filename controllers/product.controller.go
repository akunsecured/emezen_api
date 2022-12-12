package controllers

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"

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

func containsInt(s []int64, e int64) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

func (pc *ProductController) GetAllProducts(ctx *gin.Context) {
	_, err := pc.CheckHeaderAuthorization(ctx)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"message": err.Error()})
		return
	}

	nameFilter := ctx.Query("name")

	categoryQuery := ctx.Query("categories")

	categories := strings.Split(categoryQuery, ",")

	categoryFilter := []int64{}
	for _, category := range categories {
		if s, err := strconv.ParseInt(category, 10, 64); err == nil {
			categoryFilter = append(categoryFilter, s)
		}
	}

	priceFromQuery := ctx.Query("price_from")
	priceFromFilter := 0.0
	if priceFromQuery != "" {
		priceFromFilter, err = strconv.ParseFloat(priceFromQuery, 32)
		if err != nil {
			ctx.JSON(http.StatusBadGateway, gin.H{"message": err.Error()})
			return
		}
	}

	priceToQuery := ctx.Query("price_to")
	priceToFilter := 999.99
	if priceToQuery != "" {
		priceToFilter, err = strconv.ParseFloat(priceToQuery, 32)
		if err != nil {
			ctx.JSON(http.StatusBadGateway, gin.H{"message": err.Error()})
			return
		}
	}

	products, err := pc.productService.GetAllProducts()
	if err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"message": err.Error()})
		return
	}

	var results = []*models.Product{}
	for _, product := range products {
		if strings.Contains(product.Name, nameFilter) {
			if product.Price >= float32(priceFromFilter) && product.Price <= float32(priceToFilter) {
				if len(categoryFilter) == 0 {
					results = append(results, product)
				} else {
					for _, category := range categoryFilter {
						if product.Category == models.Category(category) {
							results = append(results, product)
						}
					}
				}
			}
		}
	}

	ctx.JSON(http.StatusOK, gin.H{"message": results})
}

func (pc *ProductController) UpdateProduct(ctx *gin.Context) {
	claims, err := pc.CheckHeaderAuthorization(ctx)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"message": err.Error()})
		return
	}

	productId := ctx.Param("id")
	oldProduct, err := pc.productService.GetProduct(&productId)
	if err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"message": err.Error()})
		return
	}
	userId := (*claims)["sub"].(string)

	if oldProduct.SellerID != userId {
		ctx.JSON(http.StatusUnauthorized, gin.H{"message": "tried to update other people's product"})
		return
	}

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

	product.ID = oldProduct.ID

	err = pc.productService.UpdateProduct(&product)
	if err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"message": err.Error()})
		return
	}

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

func (pc *ProductController) UploadProductImages(ctx *gin.Context) {
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

	form, err := ctx.MultipartForm()
	if err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"message": err.Error()})
		return
	}

	images := form.File["product_pictures"]
	if len(images) == 0 {
		ctx.JSON(http.StatusNotFound, gin.H{"message": "No pictures were uploaded"})
		return
	}

	fileNames := make([]string, len(images))
	for i, image := range images {
		fileName := image.Filename
		fileNames[i] = "http://localhost:8080/api/v1/product/image/" + fileName

		fmt.Println("File with name " + fileName + " is arrived.")
		filePath := "images/product_pictures/" + fileName

		file, err := image.Open()
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
			return
		}

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
	}

	product.Images = append(product.Images, fileNames...)
	pc.productService.UpdateProduct(product)

	ctx.JSON(http.StatusOK, gin.H{"message": true})
}

func (pc *ProductController) GetProductImage(ctx *gin.Context) {
	filename := ctx.Param("filename")

	ctx.File("images/product_pictures/" + filename)
}

func (pc *ProductController) GetAllProductsOfUser(ctx *gin.Context) {
	_, err := pc.CheckHeaderAuthorization(ctx)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"message": err.Error()})
		return
	}

	userId := ctx.Param("id")
	products, err := pc.productService.GetAllProductsOfUser(&userId)
	if err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"message": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": products})
}

func (pc *ProductController) BuyProducts(ctx *gin.Context) {
	claims, err := pc.CheckHeaderAuthorization(ctx)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"message": err.Error()})
		return
	}

	var cart *map[string]int32
	if err := ctx.ShouldBindJSON(&cart); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	userId := (*claims)["sub"].(string)

	err = pc.productService.BuyProducts(cart, &userId)
	if err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"message": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "Products were bought successfully"})
}

func (pc *ProductController) GetProductObserverOfUser(ctx *gin.Context) {
	claims, err := pc.CheckHeaderAuthorization(ctx)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"message": err.Error()})
		return
	}

	userId := (*claims)["sub"].(string)

	productObserver, err := pc.productService.GetProductObserverOfUser(&userId)
	if err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"message": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": productObserver})
}

func (pc *ProductController) UpdateProductObserver(ctx *gin.Context) {
	_, err := pc.CheckHeaderAuthorization(ctx)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"message": err.Error()})
		return
	}

	var productObserver *models.ProductObserver
	if err := ctx.ShouldBindJSON(&productObserver); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	updatedProductObserver, err := pc.productService.UpdateProductObserver(productObserver)
	if err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"message": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": updatedProductObserver})
}

func (pc *ProductController) RegisterProductRoutes(rg *gin.RouterGroup) {
	productRoute := rg.Group("/product")
	productRoute.POST("/create", pc.CreateProduct)
	productRoute.GET("/get/:id", pc.GetProduct)
	productRoute.GET("/get_all", pc.GetAllProducts)
	productRoute.PUT("/update/:id", pc.UpdateProduct)
	productRoute.DELETE("/delete/:id", pc.DeleteProduct)
	productRoute.POST("/image/:id", pc.UploadProductImages)
	productRoute.GET("/image/:filename", pc.GetProductImage)
	productRoute.GET("/get_all/:id", pc.GetAllProductsOfUser)
	productRoute.POST("/buy", pc.BuyProducts)
	productRoute.GET("/observer", pc.GetProductObserverOfUser)
	productRoute.PUT("/observer", pc.UpdateProductObserver)
}
