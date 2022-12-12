package services

import "github.com/akunsecured/emezen_api/models"

type ProductService interface {
	AddProduct(*models.Product) (*string, error)
	GetProduct(*string) (*models.Product, error)
	GetAllProducts() ([]*models.Product, error)
	UpdateProduct(*models.Product) error
	DeleteProduct(*string) error
	GetAllProductsOfUser(*string) ([]*models.Product, error)
	BuyProducts(*map[string]int32, *string) error
	GetProductObserverOfUser(*string) (*models.ProductObserver, error)
	UpdateProductObserver(*models.ProductObserver) (*models.ProductObserver, error)
}
