package handlers

import (
	"backend/internal/models"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"net/http"
	"strconv"
)

// HandleRemoveEventFromPublication godoc
// @Summary      Удалить событие из публикации
// @Description  Удаляет событие из черновика публикации пользователя.
// @Tags         Events-Publications
// @Accept       json
// @Produce      json
// @Param        eventID  path    int  true  "Идентификатор события"
// @Security BearerAuth
// @Success      204      {object}  nil
// @Failure      400      {object}  gin.H  "Неверный формат ID события"
// @Failure      401      {object}  gin.H  "Неверные учетные данные"
// @Failure      404      {object}  gin.H  "Событие или публикация не найдены"
// @Failure      500      {object}  gin.H  "Внутренняя ошибка сервера"
// @Router       /events/{eventID}/publications [delete]
func (h *Handler) HandleRemoveEventFromPublication(c *gin.Context) {
	userID, err := getUserIDFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	eventID, err := strconv.Atoi(c.Param("eventID"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid event ID"})
		return
	}

	err = h.repo.RemoveEventFromPublication(userID, eventID)
	if err != nil {
		if errors.Is(err, models.ErrEventNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		} else if errors.Is(err, models.ErrPublicationNotFound) {
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	c.JSON(http.StatusNoContent, nil)
}

// HandleUpdateEventPriority godoc
// @Summary      Обновить приоритет события
// @Description  Обновляет приоритет события в черновике публикации пользователя.
// @Tags         Events-Publications
// @Accept       json
// @Produce      json
// @Param        eventID  path    int                   true  "Идентификатор события"
// @Param        priority body    models.UpdateEventPriority  true  "Новый приоритет события"
// @Security BearerAuth
// @Success      200      {object}  nil
// @Failure      400      {object}  gin.H  "Неверный формат ID события или данных приоритета"
// @Failure      401      {object}  gin.H  "Неверные учетные данные"
// @Failure      404      {object}  gin.H  "Черновик публикации не найден"
// @Failure      409      {object}  gin.H  "Публикация уже удалена"
// @Failure      500      {object}  gin.H  "Внутренняя ошибка сервера"
// @Router       /events/{eventID}/priority [put]
func (h *Handler) HandleUpdateEventPriority(c *gin.Context) {
	userID, err := getUserIDFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	eventID, err := strconv.Atoi(c.Param("eventID"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid event ID"})
		return
	}

	var priority models.UpdateEventPriority
	if err = c.ShouldBindJSON(&priority); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	draftPublication, _ := h.repo.GetDraftPublication(userID)
	if draftPublication == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	err = h.repo.UpdateEventPriority(userID, draftPublication.ID, eventID, priority.Priority)
	if err != nil {
		if errors.Is(err, models.ErrPublicationNotFound) {
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, nil)
}
