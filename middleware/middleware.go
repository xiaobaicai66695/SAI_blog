package middleware

import "github.com/gin-gonic/gin"

func AutoMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := c.GetHeader("Authorization")
		if tokenString == "" {
			return
		}
		_, claims, err := ParseToken(tokenString)
		if err != nil {
			return
		}
		c.Set("uid", claims.UserId)
		c.Next()
	}
}
