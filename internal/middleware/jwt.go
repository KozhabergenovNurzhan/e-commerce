package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"ecommerce/internal/auth"
	"ecommerce/internal/models"
)

const (
	CtxUserID = "user_id"
	CtxRole   = "role"
)

func JWT(secret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		header := c.GetHeader("Authorization")
		if !strings.HasPrefix(header, "Bearer ") {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "missing token"})
			return
		}

		raw := strings.TrimPrefix(header, "Bearer ")

		claims, err := auth.ParseToken(raw, secret)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
			return
		}

		c.Set(CtxUserID, claims.UserID)
		c.Set(CtxRole, claims.Role)
		c.Next()
	}
}

func RequireRole(role string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userRole, _ := c.Get(CtxRole)
		if userRole != models.Role(role) {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "forbidden"})
			return
		}
		c.Next()
	}
}
