package services

import "github.com/akunsecured/emezen_api/models"

type ProductService interface {
	AddProduct(*models.Product) (*string, error)
	GetProduct(*string) (*models.Product, error)
	GetAllProducts() ([]*models.Product, error)
	UpdateProduct(*models.Product) error
	DeleteProduct(*string) error
}
