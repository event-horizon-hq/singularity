package middleware

import (
	"net/http"
	"singularity/internal/auth"
	"singularity/internal/enum"
	"strings"

	"github.com/gin-gonic/gin"
)

const ContextKeyClaims = "claims"

func Auth(authenticationService *auth.AuthenticationService) gin.HandlerFunc {
	return func(context *gin.Context) {
		claims, ok := extractAndValidateClaims(context, authenticationService)
		if !ok {
			return
		}

		context.Set(ContextKeyClaims, claims)
		context.Next()
	}
}

func ServerOnly(authenticationService *auth.AuthenticationService) gin.HandlerFunc {
	return func(context *gin.Context) {
		claims, ok := extractAndValidateClaims(context, authenticationService)
		if !ok {
			return
		}

		if claims.TokenType != enum.TokenTypeMaster {
			context.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"error": "This endpoint requires a server token.",
			})
			return
		}

		context.Set(ContextKeyClaims, claims)
		context.Next()
	}
}

func extractAndValidateClaims(context *gin.Context, authService *auth.AuthenticationService) (*auth.Claims, bool) {
	header := context.GetHeader("Authorization")

	if header == "" {
		context.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"error": "The header is empty!",
		})
		return nil, false
	}

	tokenString := strings.TrimPrefix(header, "Bearer ")
	claims, err := authService.Validate(tokenString)

	if err != nil {
		context.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"error": "Invalid API key.",
		})
		return nil, false
	}

	return claims, true
}

func GetClaims(context *gin.Context) *auth.Claims {
	if claims, exists := context.Get(ContextKeyClaims); exists {
		if c, ok := claims.(*auth.Claims); ok {
			return c
		}
	}
	return nil
}
