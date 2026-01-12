package middlewares

import (
	"net/http"
	"strings"

	"example.com/REST-API-Project/utils"
	"github.com/gin-gonic/gin"
)

func Authenticate(context *gin.Context) {
	authHeader := context.GetHeader("Authorization")

	if authHeader == "" {
		context.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "Not authorized."})
		return
	}

	token := strings.TrimPrefix(authHeader, "Bearer ")

	userId, err := utils.VerifyToken(token)
	if err != nil {
		context.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "Not authorized."})
		return
	}

	context.Set("userId", userId)

	context.Next()
}
