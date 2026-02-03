package auth

import (
	"singularity/internal/enum"

	"github.com/golang-jwt/jwt/v5"
)

type Claims struct {
	jwt.RegisteredClaims
	TokenType enum.TokenType `json:"type"`
}
