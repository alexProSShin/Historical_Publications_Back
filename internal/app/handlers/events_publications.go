package handlers

import (
	"backend/internal/lib/api/resp"
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
// @Success      204
// @Failure      400      {object}  resp.ErrorResponse   "Неверный формат ID события"
// @Failure      401      {object}  resp.ErrorResponse   "Неверные учетные данные"
// @Failure      404      {object}  resp.ErrorResponse   "Событие или публикация не найдены"
// @Failure      500      {object}  resp.ErrorResponse   "Внутренняя ошибка сервера"
// @Router       /events/{eventID}/publications [delete]
func (h *Handler) HandleRemoveEventFromPublication(c *gin.Context) {
	userID, err := getUserIDFromContext(c)
	if err != nil {
		resp.WriteError(c.Writer, http.StatusUnauthorized, resp.ErrorDetailList("error", err.Error()), nil)
		return
	}

	eventID, err := strconv.Atoi(c.Param("eventID"))
	if err != nil {
		resp.WriteError(c.Writer, http.StatusBadRequest, resp.ErrorDetailList("eventID", "invalid event ID"), nil)
		return
	}

	err = h.repo.RemoveEventFromPublication(userID, eventID)
	if err != nil {
		if errors.Is(err, models.ErrEventNotFound) {
			resp.WriteError(c.Writer, http.StatusNotFound, resp.SingleError(err.Error()), nil)
		} else if errors.Is(err, models.ErrPublicationNotFound) {
			resp.WriteError(c.Writer, http.StatusConflict, resp.SingleError(err.Error()), nil)
		} else {
			resp.WriteError(c.Writer, http.StatusInternalServerError, resp.SingleError(err.Error()), nil)
		}
		return
	}

	resp.WriteJSON(c.Writer, http.StatusNoContent, nil)
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
// @Success      200
// @Failure      400      {object}  resp.ErrorResponse   "Неверный формат ID события или данных приоритета"
// @Failure      401      {object}  resp.ErrorResponse   "Неверные учетные данные"
// @Failure      404      {object}  resp.ErrorResponse   "Черновик публикации не найден"
// @Failure      409      {object}  resp.ErrorResponse   "Публикация уже удалена"
// @Failure      500      {object}  resp.ErrorResponse   "Внутренняя ошибка сервера"
// @Router       /events/{eventID}/priority [put]
func (h *Handler) HandleUpdateEventPriority(c *gin.Context) {
	userID, err := getUserIDFromContext(c)
	if err != nil {
		resp.WriteError(c.Writer, http.StatusUnauthorized, resp.ErrorDetailList("error", err.Error()), nil)
		return
	}

	eventID, err := strconv.Atoi(c.Param("eventID"))
	if err != nil {
		resp.WriteError(c.Writer, http.StatusBadRequest, resp.ErrorDetailList("eventID", "invalid event ID"), nil)
		return
	}

	var priority models.UpdateEventPriority
	if err = c.ShouldBindJSON(&priority); err != nil {
		resp.WriteError(c.Writer, http.StatusBadRequest, resp.ErrorDetailList("request_body", err.Error()), nil)
		return
	}

	draftPublication, _ := h.repo.GetDraftPublication(userID)
	if draftPublication == nil {
		resp.WriteError(c.Writer, http.StatusNotFound, resp.SingleError("draft publication not found"), nil)
		return
	}

	err = h.repo.UpdateEventPriority(userID, draftPublication.ID, eventID, priority.Priority)
	if err != nil {
		if errors.Is(err, models.ErrPublicationNotFound) {
			resp.WriteError(c.Writer, http.StatusConflict, resp.SingleError(err.Error()), nil)
		} else {
			resp.WriteError(c.Writer, http.StatusInternalServerError, resp.SingleError(err.Error()), nil)
		}
		return
	}

	resp.WriteJSON(c.Writer, http.StatusOK, nil)
}
