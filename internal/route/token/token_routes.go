package token

import (
	"net/http"
	"singularity/internal/auth"
	"singularity/internal/enum"

	"github.com/gin-gonic/gin"
)

func Register(group *gin.RouterGroup, authService *auth.AuthenticationService) {
	group.POST("", func(c *gin.Context) {
		GenerateUserTokenHandler(c, authService)
	})
}

func GenerateUserTokenHandler(c *gin.Context, authService *auth.AuthenticationService) {
	token, err := authService.GenerateToken(enum.TokenTypeSlave)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to generate token",
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"token": token,
	})
}
