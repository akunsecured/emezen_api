package security

import (
	"fmt"
	"github.com/akunsecured/emezen_api/models"
	"github.com/akunsecured/emezen_api/utils"
	"github.com/form3tech-oss/jwt-go"
	"strings"
	"time"
)

var (
	JwtSecretKey = []byte("Emezen_SUP3R_S3CR3T_K3Y")
)

type JwtUserClaims struct {
	User models.User `json:"user"`
	jwt.StandardClaims
}

func NewAccessToken(user models.User) (string, error) {
	userId := user.ID.Hex()
	claims := JwtUserClaims{
		user,
		jwt.StandardClaims{
			Subject:   userId,
			IssuedAt:  time.Now().Unix(),
			ExpiresAt: time.Now().Add(time.Hour * 1).Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(JwtSecretKey)
}

func NewRefreshToken(userId string) (string, error) {
	claims := jwt.StandardClaims{
		Subject:   userId,
		IssuedAt:  time.Now().Unix(),
		ExpiresAt: time.Now().Add(time.Hour * 48).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(JwtSecretKey)
}

func CreateAccessAndRefreshTokens(user models.User) (*models.WrappedToken, error) {
	accessToken, err := NewAccessToken(user)
	if err != nil {
		return nil, err
	}

	refreshToken, err := NewRefreshToken(user.ID.Hex())
	if err != nil {
		return nil, err
	}

	return &models.WrappedToken{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func ValidateSignedMethod(token *jwt.Token) (interface{}, error) {
	if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
		return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
	}
	return JwtSecretKey, nil
}

func ParseToken(tokenString string) (*jwt.MapClaims, error) {
	if !strings.HasPrefix(tokenString, "Bearer ") {
		return nil, utils.ErrInvalidTokenFormat
	}
	tokenString = strings.Split(tokenString, " ")[1]
	token, err := jwt.Parse(tokenString, ValidateSignedMethod)
	if err != nil {
		return nil, err
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, utils.ErrTokenParseError
	}

	return &claims, nil
}
