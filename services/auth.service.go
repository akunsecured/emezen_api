package services

import (
	"github.com/akunsecured/emezen_api/models"
	"github.com/form3tech-oss/jwt-go"
)

type AuthService interface {
	Register(*models.UserDataWithCredentials) (*models.WrappedToken, error)
	Login(*models.UserCredentials) (*models.WrappedToken, error)
	Update(*models.UserCredentials) error
	NewAccessToken(*jwt.MapClaims) (*string, error)
	CurrentUser(*jwt.MapClaims) (*models.User, error)
	DeleteUser(*jwt.MapClaims) error
}
