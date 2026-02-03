package auth

import (
	"fmt"
	"time"

	"singularity/internal/enum"

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

func (authService *AuthenticationService) GenerateToken(tokenType enum.TokenType) (string, error) {
	claims := Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Now().AddDate(10, 0, 0)),
		},
		TokenType: tokenType,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(authService.SecretKey))
}

func (authService *AuthenticationService) GenerateSystemToken() (string, error) {
	return authService.GenerateToken(enum.TokenTypeMaster)
}

func (authService *AuthenticationService) Validate(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(authService.SecretKey), nil
	})

	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, fmt.Errorf("invalid token claims")
	}

	if claims.TokenType == "" {
		claims.TokenType = enum.TokenTypeMaster
	}

	return claims, nil
}
