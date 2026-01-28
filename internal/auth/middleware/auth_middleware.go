package middleware

import (
	"net/http"
	"singularity/internal/auth"
	"strings"

	"github.com/gin-gonic/gin"
)

func Auth(authenticationService *auth.AuthenticationService) gin.HandlerFunc {
	return func(context *gin.Context) {
		header := context.GetHeader("Authorization")

		if header == "" {
			context.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"message": "The header is empty!",
			})

			return
		}

		tokenString := strings.TrimPrefix(header, "Bearer ")
		token, err := authenticationService.Validate(tokenString)

		if err != nil || !token.Valid {
			context.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"message": "Invalid API key.",
			})

			return
		}

		context.Next()
	}
}
