package handlers

import (
	"backend/internal/app/repository"
	"backend/internal/models"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/pkg/errors"
	"net/http"
	"time"
)

func getUserIDFromContext(c *gin.Context) (int, error) {
	userID, exists := c.Get("userID")
	if !exists {
		return 0, errors.New("User ID not found in context")
	}

	return userID.(int), nil
}

var jwtSecretKey = []byte("your_secret_key")

func generateJWT(userID int) (string, error) {
	claims := jwt.MapClaims{
		"userID": userID,
		"exp":    time.Now().Add(24 * time.Hour).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecretKey)
}

func saveTokenToRedis(redisRepo repository.RedisRepo, token string, userID int) error {
	return redisRepo.Set("user:"+string(userID), token, time.Hour*24)
}

func validateTokenInRedis(redisRepo repository.RedisRepo, userID int) (string, error) {
	token, err := redisRepo.Get("user:" + string(userID))
	if err != nil {
		return "", err
	}
	return token, nil
}

func (h *Handler) HandleRegisterUser(c *gin.Context) {
	var user models.RegisterUserDTO
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request data"})
		return
	}

	newUser, err := h.repo.CreateUser(&user)
	if err != nil {
		if errors.Is(err, models.ErrUserAlreadyExists) {
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	token, err := generateJWT(newUser.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if err = saveTokenToRedis(h.redis, token, newUser.ID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"token": token, "user": newUser})
}

func (h *Handler) HandleLoginUser(c *gin.Context) {
	var credentials models.LoginUserDTO
	if err := c.ShouldBindJSON(&credentials); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request data"})
		return
	}

	user, err := h.repo.AuthenticateUser(credentials.Email, credentials.Password)
	if err != nil {
		if errors.Is(err, models.ErrInvalidCredentials) {
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	token, err := generateJWT(user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if err = saveTokenToRedis(h.redis, token, user.ID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token, "user": user})
}

func (h *Handler) HandleLogoutUser(c *gin.Context) {
	userID, err := getUserIDFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	if err = h.redis.Delete("user:" + string(userID)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusNoContent, nil)
}

func (h *Handler) HandleUpdateUser(c *gin.Context) {
	userID, err := getUserIDFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	var updateData models.UpdateUserDTO
	if err := c.ShouldBindJSON(&updateData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request data"})
		return
	}

	updatedUser, err := h.repo.UpdateUser(userID, &updateData)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, updatedUser)
}
