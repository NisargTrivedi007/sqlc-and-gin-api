package middleware

import (
	"net/http"
	"sqlc_api/jwt_code"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {

		tokenString := c.GetHeader("Authorization")
		if tokenString == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Token required"})
			c.Abort()
			return
		}

		token, err := jwt_code.ValidateToken(tokenString)
		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		// Extract claims and store in context
		claims, ok := token.Claims.(jwt.MapClaims)
		if ok && token.Valid {
			c.Set("username", claims["username"])
		}

		c.Next()
	}
}
