package auth

import (
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func Middleware(as AuthService) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			c.JSON(401, gin.H{"error": "unauthorized"})
			c.Abort()
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		token, err := jwt.ParseWithClaims(tokenString, &CourtlyClaims{}, func(token *jwt.Token) (interface{}, error) {
			if token.Method != jwt.SigningMethodHS256 {
				return nil, jwt.ErrTokenSignatureInvalid
			}

			return as.GetSecretKey(), nil
		})

		if err != nil || !token.Valid {
			c.JSON(401, gin.H{"error": "unauthorized"})
			c.Abort()
			return
		}

		claims, ok := token.Claims.(*CourtlyClaims)
		if !ok {
			c.JSON(401, gin.H{"error": "unauthorized"})
			c.Abort()
			return
		}

        if claims.ExpiresAt.Time.Before(time.Now()) {
			c.JSON(401, gin.H{"error": "token expired"})
			c.Abort(); return
		}

		c.Set("company_id", claims.Sub)
		c.Set("jwt_token", tokenString)

		c.Next()
	}
}
