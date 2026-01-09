package auth

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type AuthenticationService struct {
	SecretKey string
}

func NewAuthenticationService(secretKey string) *AuthenticationService {
	return &AuthenticationService{
		SecretKey: secretKey,
	}
}

func (authService *AuthenticationService) GenerateSystemToken() (string, error) {
	claims := jwt.MapClaims{
		"iat": time.Now().Unix(),
		"exp": time.Now().AddDate(10, 0, 0).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(authService.SecretKey))
}

func (authService *AuthenticationService) Validate(tokenString string) (*jwt.Token, error) {
	return jwt.Parse(tokenString, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(authService.SecretKey), nil
	})
}
