package middleware

import (
	"backend/internal/app/repository"
	"backend/internal/lib/api/resp"
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
			resp.WriteError(c.Writer, http.StatusUnauthorized, resp.SingleError("authorization header is required"), nil)
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
			resp.WriteError(c.Writer, http.StatusUnauthorized, resp.SingleError("invalid or expired token"), nil)
			c.Abort()
			return
		}

		redisToken, err := redisRepo.Get("user:" + strconv.Itoa(userID))
		if err != nil || redisToken != tokenString {
			resp.WriteError(c.Writer, http.StatusUnauthorized, resp.SingleError("token is invalid or revoked"), nil)
			c.Abort()
			return
		}

		c.Set("userID", userID)
		c.Next()
	}
}

func OptionalAuthMiddleware(redisRepo repository.RedisRepo) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")

		if authHeader == "" {
			c.Next()
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
			c.Next()
			return
		}

		redisToken, err := redisRepo.Get("user:" + strconv.Itoa(userID))
		if err != nil || redisToken != tokenString {
			c.Next()
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
			resp.WriteError(c.Writer, http.StatusUnauthorized, resp.SingleError("user ID not found in context"), nil)
			c.Abort()
			return
		}

		user, err := getUser(userID.(int))
		if err != nil || user.Role != models.RoleModerator {
			resp.WriteError(c.Writer, http.StatusForbidden, resp.SingleError("access denied: Moderator role required"), nil)
			c.Abort()
			return
		}

		c.Next()
	}
}
