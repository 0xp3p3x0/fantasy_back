package middleware

import (
	"net/http"
	"strings"

	"back/internal/model"
	"back/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
)

func TokenAuth(expectedToken string) gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("X-API-Token")
		if token == "" {
			authHeader := c.GetHeader("Authorization")
			token = strings.TrimPrefix(authHeader, "Bearer ")
		}

		if token == "" || token != expectedToken {
			c.AbortWithStatusJSON(http.StatusUnauthorized, model.APIResponse{
				Success: false,
				Message: "unauthorized",
			})
			return
		}

		c.Next()
	}
}

func JWTAuth(secret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, model.APIResponse{
				Success: false,
				Message: "missing authorization header",
			})
			return
		}

		token := strings.TrimPrefix(authHeader, "Bearer ")
		if token == authHeader || token == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, model.APIResponse{
				Success: false,
				Message: "invalid authorization format",
			})
			return
		}

		claims, err := service.ValidateJWT(token, secret)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, model.APIResponse{
				Success: false,
				Message: "invalid token",
			})
			return
		}

		c.Set("claims", claims)
		c.Next()
	}
}

func RequireRole(roles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		claimsAny, ok := c.Get("claims")
		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, model.APIResponse{
				Success: false,
				Message: "missing auth claims",
			})
			return
		}

		claims, ok := claimsAny.(*jwt.MapClaims)
		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, model.APIResponse{
				Success: false,
				Message: "invalid auth claims",
			})
			return
		}

		role, _ := (*claims)["role"].(string)
		for _, allowed := range roles {
			if role == allowed {
				c.Next()
				return
			}
		}

		c.AbortWithStatusJSON(http.StatusForbidden, model.APIResponse{
			Success: false,
			Message: "forbidden",
		})
	}
}
