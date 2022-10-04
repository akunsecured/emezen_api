package services

import (
	"github.com/akunsecured/emezen_api/models"
)

type UserService interface {
	CreateUser(*models.User) (*string, error)
	GetUser(*string) (*models.User, error)
	UpdateUser(*models.User) error
	DeleteUser(*string) error
}
