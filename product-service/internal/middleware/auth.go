package middleware

import (
	"net/http"

	"github.com/ecommerce/product-service/internal/utils"
	"github.com/gin-gonic/gin"
)

func AuthMiddleware(jwtSecret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		token, err := utils.ExtractToken(authHeader)
		if err != nil {
			utils.ErrorResponse(c, http.StatusUnauthorized, err.Error(), "UNAUTHORIZED", nil)
			c.Abort()
			return
		}

		claims, err := utils.ValidateToken(token, jwtSecret)
		if err != nil {
			utils.ErrorResponse(c, http.StatusUnauthorized, "Invalid or expired token", "INVALID_TOKEN", err.Error())
			c.Abort()
			return
		}

		// Set user info in context
		c.Set("userID", claims.UserID)
		c.Set("email", claims.Email)
		c.Set("role", claims.Role)

		c.Next()
	}
}

func RoleMiddleware(allowedRoles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		role, exists := c.Get("role")
		if !exists {
			utils.ErrorResponse(c, http.StatusUnauthorized, "User role not found", "UNAUTHORIZED", nil)
			c.Abort()
			return
		}

		roleStr := role.(string)
		allowed := false
		for _, allowedRole := range allowedRoles {
			if roleStr == allowedRole {
				allowed = true
				break
			}
		}

		if !allowed {
			utils.ErrorResponse(c, http.StatusForbidden, "Insufficient permissions", "FORBIDDEN", nil)
			c.Abort()
			return
		}

		c.Next()
	}
}
