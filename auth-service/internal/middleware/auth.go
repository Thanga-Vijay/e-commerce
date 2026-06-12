package middleware

import (
	"net/http"
	"strings"

	"github.com/ecommerce/auth-service/internal/utils"
	"github.com/gin-gonic/gin"
)

func AuthMiddleware(jwtManager *utils.JWTManager) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			utils.ErrorResponse(c, http.StatusUnauthorized, "Authorization header required", "UNAUTHORIZED", nil)
			c.Abort()
			return
		}

		// Extract token from "Bearer <token>"
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			utils.ErrorResponse(c, http.StatusUnauthorized, "Invalid authorization header format", "INVALID_AUTH_HEADER", nil)
			c.Abort()
			return
		}

		token := parts[1]
		claims, err := jwtManager.ValidateToken(token)
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
