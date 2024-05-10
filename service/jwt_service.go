package service

import (
	"time"
	"vcs-sms/model/entity"

	"github.com/golang-jwt/jwt"
)

type JWTService struct {
}

func NewJWTService() *JWTService {
	return &JWTService{}
}

func (service *JWTService) GenerateToken(user *entity.User) (string, error) {
	var scopes []string
	for _, scope := range user.Scopes {
		scopes = append(scopes, scope.Name)
	}
	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": user.Username,
		"exp":      time.Now().Add(time.Hour * 24).Unix(),
		"iat":      time.Now().Unix(),
		"scopes":   scopes,
	})
	token, err := claims.SignedString([]byte("secret"))
	if err != nil {
		return "", err
	}
	return token, err
}
