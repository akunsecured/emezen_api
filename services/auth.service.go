package services

import "github.com/akunsecured/emezen_api/models"

type AuthService interface {
	Register(*models.UserCredentials) error
	Login(*models.UserCredentials) (*string, error)
	Update(*models.UserCredentials) error
}
