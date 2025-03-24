package middleware

import (
	"net/http"
	"strings"

	"qckim-backend/internal/logger"
	"qckim-backend/utils/jwt"

	"github.com/gin-gonic/gin"
)

type UserContextKey string

const UserKey UserContextKey = "userID"

var FreeRoutes = map[string]struct{}{
	"/api/users/register":               {},
	"/api/users/activate":               {},
	"/api/users/login":                  {},
	"/api/users/request-password-reset": {},
	"/api/users/reset-password":         {},
	"/api/qck/shared-data-link":         {},
	"/api/users/refresh-token":          {},
}

// JWTAuth() - handler запросов с проверкой token
func JWTAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		if _, ok := FreeRoutes[c.Request.URL.Path]; ok {
			c.Next()
			return
		}

		tokenString := c.GetHeader("Authorization")
		if tokenString == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization token is missing"})
			c.Abort()
			return
		}

		tokenString = strings.TrimPrefix(tokenString, "Bearer ")

		claims, err := jwt.ValidateJWT(tokenString)
		if err != nil {
			logger.Error("jwt.ValidateJWT(tokenString)", "err", err)
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
			c.Abort()
			return
		}

		c.Set(string(UserKey), claims.UserID)

		c.Next()
	}
}
