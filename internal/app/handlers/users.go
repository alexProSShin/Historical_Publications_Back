package handlers

import (
	"backend/internal/lib/api/resp"
	"backend/internal/models"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/pkg/errors"
	"net/http"
	"strconv"
	"time"
)

func getUserIDFromContext(c *gin.Context) (int, error) {
	userID, exists := c.Get("userID")
	if !exists {
		return 0, models.ErrUserIDNotFound
	}

	return userID.(int), nil
}

var jwtSecretKey = []byte("secret")

func generateJWT(userID int) (string, error) {
	claims := jwt.MapClaims{
		"userID": userID,
		"exp":    time.Now().Add(24 * time.Hour).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecretKey)
}

func (h *Handler) saveTokenToRedis(token string, userID int) error {
	return h.redis.Set("user:"+strconv.Itoa(userID), token, time.Hour*24)
}

// HandleRegisterUser godoc
// @Summary      Регистрация нового пользователя
// @Description  Регистрирует нового пользователя и возвращает JWT токен.
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param        user  body    models.RegisterUserDTO  true  "Данные пользователя"
// @Success      201  {object}  models.LoginResponseDTO
// @Failure      400  {object}  resp.ErrorResponse   "Неверные данные запроса"
// @Failure      409  {object}  resp.ErrorResponse   "Пользователь уже существует"
// @Failure      500  {object}  resp.ErrorResponse   "Внутренняя ошибка сервера"
// @Router       /users/register [post]
func (h *Handler) HandleRegisterUser(c *gin.Context) {
	var user models.RegisterUserDTO
	if err := c.ShouldBindJSON(&user); err != nil {
		resp.WriteError(c.Writer, http.StatusBadRequest, resp.ErrorDetailList("request_body", err.Error()), nil)
		return
	}

	newUser, err := h.repo.CreateUser(&user)
	if err != nil {
		if errors.Is(err, models.ErrUserAlreadyExists) {
			resp.WriteError(c.Writer, http.StatusConflict, resp.SingleError(err.Error()), nil)
		} else {
			resp.WriteError(c.Writer, http.StatusInternalServerError, resp.SingleError(err.Error()), nil)
		}
		return
	}

	token, err := generateJWT(newUser.ID)
	if err != nil {
		resp.WriteError(c.Writer, http.StatusInternalServerError, resp.SingleError(err.Error()), nil)
		return
	}

	if err = h.saveTokenToRedis(token, newUser.ID); err != nil {
		resp.WriteError(c.Writer, http.StatusInternalServerError, resp.SingleError(err.Error()), nil)
		return
	}

	resp.WriteJSON(c.Writer, http.StatusCreated, models.LoginResponseDTO{
		Token: token,
		User:  *newUser,
	})
}

// HandleLoginUser godoc
// @Summary      Вход пользователя
// @Description  Осуществляет вход пользователя, если его учетные данные верны, и возвращает JWT токен.
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param        credentials  body    models.LoginUserDTO  true  "Данные для входа"
// @Success      200          {object}  models.LoginResponseDTO
// @Failure      400          {object}  resp.ErrorResponse   "Неверные данные запроса"
// @Failure      401          {object}  resp.ErrorResponse   "Неверные учетные данные"
// @Failure      500          {object}  resp.ErrorResponse   "Внутренняя ошибка сервера"
// @Router       /users/login [post]
func (h *Handler) HandleLoginUser(c *gin.Context) {
	var credentials models.LoginUserDTO
	if err := c.ShouldBindJSON(&credentials); err != nil {
		resp.WriteError(c.Writer, http.StatusBadRequest, resp.ErrorDetailList("request_body", err.Error()), nil)
		return
	}

	user, err := h.repo.AuthenticateUser(credentials.Email, credentials.Password)
	if err != nil {
		if errors.Is(err, models.ErrInvalidCredentials) {
			resp.WriteError(c.Writer, http.StatusUnauthorized, resp.SingleError(err.Error()), nil)
		} else {
			resp.WriteError(c.Writer, http.StatusInternalServerError, resp.SingleError(err.Error()), nil)
		}
		return
	}

	token, err := generateJWT(user.ID)
	if err != nil {
		resp.WriteError(c.Writer, http.StatusInternalServerError, resp.SingleError(err.Error()), nil)
		return
	}

	if err = h.saveTokenToRedis(token, user.ID); err != nil {
		resp.WriteError(c.Writer, http.StatusInternalServerError, resp.SingleError(err.Error()), nil)
		return
	}

	resp.WriteJSON(c.Writer, http.StatusCreated, models.LoginResponseDTO{
		Token: token,
		User:  *user,
	})
}

// HandleLogoutUser godoc
// @Summary      Выход пользователя
// @Description  Завершаем сессию пользователя и удаляем его токен.
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Security BearerAuth
// @Success      204
// @Failure      401  {object}  resp.ErrorResponse   "Неавторизованный пользователь"
// @Failure      500  {object}  resp.ErrorResponse   "Внутренняя ошибка сервера"
// @Router       /users/logout [post]
func (h *Handler) HandleLogoutUser(c *gin.Context) {
	userID, err := getUserIDFromContext(c)
	if err != nil {
		resp.WriteError(c.Writer, http.StatusUnauthorized, resp.SingleError(err.Error()), nil)
		return
	}

	if err = h.redis.Delete("user:" + strconv.Itoa(userID)); err != nil {
		resp.WriteError(c.Writer, http.StatusInternalServerError, resp.SingleError(err.Error()), nil)
		return
	}

	resp.WriteJSON(c.Writer, http.StatusNoContent, nil)
}

// HandleUpdateUser godoc
// @Summary      Обновить информацию о пользователе
// @Description  Обновляет данные пользователя.
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param        updateData  body    models.UpdateUserDTO  true  "Данные для обновления"
// @Security BearerAuth
// @Success      200         {object}  models.User  "Обновленный пользователь"
// @Failure      400         {object}  resp.ErrorResponse   "Неверные данные запроса"
// @Failure      401         {object}  resp.ErrorResponse   "Неавторизованный пользователь"
// @Failure      500         {object}  resp.ErrorResponse   "Внутренняя ошибка сервера"
// @Router       /users/me [put]
func (h *Handler) HandleUpdateUser(c *gin.Context) {
	userID, err := getUserIDFromContext(c)
	if err != nil {
		resp.WriteError(c.Writer, http.StatusUnauthorized, resp.SingleError(err.Error()), nil)
		return
	}

	var updateData models.UpdateUserDTO
	if err := c.ShouldBindJSON(&updateData); err != nil {
		resp.WriteError(c.Writer, http.StatusBadRequest, resp.ErrorDetailList("request_body", err.Error()), nil)
		return
	}

	updatedUser, err := h.repo.UpdateUser(userID, &updateData)
	if err != nil {
		resp.WriteError(c.Writer, http.StatusInternalServerError, resp.SingleError(err.Error()), nil)
		return
	}

	resp.WriteJSON(c.Writer, http.StatusOK, updatedUser)
}

func (h *Handler) GetUserByID(userID int) (*models.User, error) {
	return h.repo.GetUserByID(userID)
}
