package services

import (
	"github.com/akunsecured/emezen_api/models"
)

type AuthService interface {
	Register(*models.UserDataWithCredentials) (*models.WrappedToken, error)
	Login(*models.UserCredentials) (*models.WrappedToken, error)
	Update(*models.UserCredentials) error
	NewAccessToken(string) (*string, error)
}
