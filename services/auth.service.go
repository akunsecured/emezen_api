package services

import "github.com/akunsecured/emezen_api/models"

type AuthService interface {
	Register(*models.UserDataWithCredentials) (*string, error)
	Login(*models.UserCredentials) (*string, error)
	Update(*models.UserCredentials) error
}
