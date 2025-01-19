package middleware

import "github.com/gin-gonic/gin"

func InjectUserIDMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := 1
		c.Set("userID", userID)
		c.Next()
	}
}
