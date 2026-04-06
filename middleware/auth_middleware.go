package middleware

import (
	"chattrix/utils"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware(jwtService *utils.JWTService) gin.HandlerFunc {
	return func(context *gin.Context) {
		authHeader := context.GetHeader("Authorization")
		if authHeader == "" {
			utils.ErrorResponse(context, http.StatusUnauthorized, "authorization header missing")
			context.Abort()
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			utils.ErrorResponse(context, http.StatusUnauthorized, "invalid authorization header format")
			context.Abort()
			return
		}
		
    token := parts[1]

		claims, err := jwtService.ValidateToken(token)
    if err != nil {
      utils.ErrorResponse(context, http.StatusUnauthorized, "invalid or expired token")
      context.Abort()
      return
    }

		context.Set("user_id", claims.UserID)
		context.Set("email", claims.Username)

		context.Next()
	}
}