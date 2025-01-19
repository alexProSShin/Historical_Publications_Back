package middleware

import (
	"backend/internal/app/repository"
	"backend/internal/models"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/pkg/errors"
	"net/http"
	"strconv"
	"strings"
)

var jwtSecretKey = []byte("secret")

func parseJWT(tokenString string) (int, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrSignatureInvalid
		}
		return jwtSecretKey, nil
	})

	if err != nil || !token.Valid {
		return 0, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return 0, errors.New("invalid token claims")
	}

	userID, ok := claims["userID"].(float64)
	if !ok {
		return 0, errors.New("invalid token payload")
	}

	return int(userID), nil
}

func AuthMiddleware(redisRepo repository.RedisRepo) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header is required"})
			c.Abort()
			return
		}

		var tokenString string
		if strings.HasPrefix(authHeader, "Bearer ") {
			tokenString = strings.TrimPrefix(authHeader, "Bearer ")
		} else {
			tokenString = authHeader
		}

		userID, err := parseJWT(tokenString)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid or expired token"})
			c.Abort()
			return
		}

		redisToken, err := redisRepo.Get("user:" + strconv.Itoa(userID))
		if err != nil || redisToken != tokenString {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "token is invalid or revoked"})
			c.Abort()
			return
		}

		c.Set("userID", userID)
		c.Next()
	}
}

func RequireModeratorRole(getUser func(int) (*models.User, error)) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, exists := c.Get("userID")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "User ID not found in context"})
			c.Abort()
			return
		}

		user, err := getUser(userID.(int))
		if err != nil || user.Role != models.RoleModerator {
			c.JSON(http.StatusForbidden, gin.H{"error": "Access denied: Moderator role required"})
			c.Abort()
			return
		}

		c.Next()
	}
}
