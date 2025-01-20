package handlers

import (
	"backend/internal/lib/api/resp"
	"backend/internal/models"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"net/http"
	"strconv"
	"time"
)

// HandleGetPublications godoc
// @Summary      Получить список публикаций
// @Description  Возвращает список публикаций с возможностью фильтрации по статусу и диапазону дат.
// @Tags         Publications
// @Accept       json
// @Produce      json
// @Param        status    query    string     false  "Статус публикации"
// @Param        startDate query    string     false  "Дата начала (формат: yyyy-mm-dd)"
// @Param        endDate   query    string     false  "Дата окончания (формат: yyyy-mm-dd)"
// @Security BearerAuth
// @Success      200      {array}   models.Publication   "Список публикаций"
// @Failure      400      {object}  resp.ErrorResponse  "Неверный формат параметров"
// @Failure      401      {object}  resp.ErrorResponse  "Неверные учетные данные"
// @Failure      404      {object}  resp.ErrorResponse  "Публикации не найдены"
// @Failure      500      {object}  resp.ErrorResponse  "Внутренняя ошибка сервера"
// @Router       /publications [get]
func (h *Handler) HandleGetPublications(c *gin.Context) {
	userID, err := getUserIDFromContext(c)
	if err != nil {
		resp.WriteError(c.Writer, http.StatusUnauthorized, resp.SingleError(err.Error()), nil)
		return
	}

	status := models.PublicationStatus(c.Query("status"))

	startDateStr := c.Query("startDate")
	endDateStr := c.Query("endDate")

	var startDate, endDate *time.Time
	if startDateStr != "" {
		parsedStartDate, err := time.Parse("2006-01-02", startDateStr)
		if err != nil {
			resp.WriteError(c.Writer, http.StatusBadRequest, resp.ErrorDetailList("startDate", err.Error()), nil)
			return
		}
		startDate = &parsedStartDate
	}

	if endDateStr != "" {
		parsedEndDate, err := time.Parse("2006-01-02", endDateStr)
		if err != nil {
			resp.WriteError(c.Writer, http.StatusBadRequest, resp.ErrorDetailList("endDate", err.Error()), nil)
			return
		}
		endDate = &parsedEndDate
	}

	publications, err := h.repo.GetPublications(userID, status, startDate, endDate)
	if err != nil {
		resp.WriteError(c.Writer, http.StatusInternalServerError, resp.SingleError(err.Error()), nil)
		return
	}

	if len(publications) == 0 {
		resp.WriteError(c.Writer, http.StatusNotFound, resp.SingleError(err.Error()), nil)
		return
	}

	resp.WriteJSON(c.Writer, http.StatusOK, publications)
}

// HandleGetPublicationByID godoc
// @Summary      Получить публикацию по ID
// @Description  Возвращает публикацию по уникальному идентификатору.
// @Tags         Publications
// @Accept       json
// @Produce      json
// @Param        publicationID  path    int  true  "Идентификатор публикации"
// @Security BearerAuth
// @Success      200            {object}  models.Publication   "Публикация"
// @Failure      400            {object}  resp.ErrorResponse  "Неверный формат ID публикации"
// @Failure      401            {object}  resp.ErrorResponse  "Неверные учетные данные"
// @Failure      404            {object}  resp.ErrorResponse  "Публикация не найдена"
// @Failure      500            {object}  resp.ErrorResponse  "Внутренняя ошибка сервера"
// @Router       /publications/{publicationID} [get]
func (h *Handler) HandleGetPublicationByID(c *gin.Context) {
	userID, err := getUserIDFromContext(c)
	if err != nil {
		resp.WriteError(c.Writer, http.StatusUnauthorized, resp.SingleError(err.Error()), nil)
		return
	}

	id, err := strconv.Atoi(c.Param("publicationID"))
	if err != nil {
		resp.WriteError(c.Writer, http.StatusBadRequest, resp.ErrorDetailList("publicationID", err.Error()), nil)
		return
	}

	publication, err := h.repo.GetPublicationByID(userID, id)
	if err != nil {
		if errors.Is(err, models.ErrPublicationNotFound) {
			resp.WriteError(c.Writer, http.StatusNotFound, resp.SingleError(err.Error()), nil)
		} else {
			resp.WriteError(c.Writer, http.StatusInternalServerError, resp.SingleError(err.Error()), nil)
		}
		return
	}

	resp.WriteJSON(c.Writer, http.StatusOK, publication)
}

// HandleUpdatePublication godoc
// @Summary      Обновить публикацию по ID
// @Description  Обновляет информацию о публикации по уникальному идентификатору.
// @Tags         Publications
// @Accept       json
// @Produce      json
// @Param        publicationID  path    int                        true  "Идентификатор публикации"
// @Param        publication    body    models.UpdatePublicationDTO true  "Обновленные данные публикации"
// @Security BearerAuth
// @Success      200            {object}  models.Publication   "Обновленная публикация"
// @Failure      400            {object}  resp.ErrorResponse  "Неверный формат данных"
// @Failure      401            {object}  resp.ErrorResponse  "Неверные учетные данные"
// @Failure      404            {object}  resp.ErrorResponse  "Публикация не найдена"
// @Failure      500            {object}  resp.ErrorResponse  "Внутренняя ошибка сервера"
// @Router       /publications/{publicationID} [put]
func (h *Handler) HandleUpdatePublication(c *gin.Context) {
	userID, err := getUserIDFromContext(c)
	if err != nil {
		resp.WriteError(c.Writer, http.StatusUnauthorized, resp.SingleError(err.Error()), nil)
		return
	}

	id, err := strconv.Atoi(c.Param("publicationID"))
	if err != nil {
		resp.WriteError(c.Writer, http.StatusBadRequest, resp.ErrorDetailList("publicationID", err.Error()), nil)
		return
	}

	var publication models.UpdatePublicationDTO
	if err := c.ShouldBindJSON(&publication); err != nil {
		resp.WriteError(c.Writer, http.StatusBadRequest, resp.ErrorDetailList("request_body", err.Error()), nil)
		return
	}

	updatedPublication, err := h.repo.UpdatePublication(userID, id, &publication)
	if err != nil {
		if errors.Is(err, models.ErrPublicationNotFound) {
			resp.WriteError(c.Writer, http.StatusNotFound, resp.SingleError(err.Error()), nil)
		} else {
			resp.WriteError(c.Writer, http.StatusInternalServerError, resp.SingleError(err.Error()), nil)
		}
		return
	}

	resp.WriteJSON(c.Writer, http.StatusOK, updatedPublication)
}

// HandleFormPublication godoc
// @Summary      Перевести публикацию в статус работы
// @Description  Переводит черновик публикации в статус работы.
// @Tags         Publications
// @Accept       json
// @Produce      json
// @Param        publicationID  path    int  true  "Идентификатор публикации"
// @Security BearerAuth
// @Success      200            {object}  models.Publication   "Публикация в статусе работы"
// @Failure      400            {object}  resp.ErrorResponse  "Неверный формат ID публикации"
// @Failure      401            {object}  resp.ErrorResponse  "Неверные учетные данные"
// @Failure      403            {object}  resp.ErrorResponse  "Недостаточно прав"
// @Failure      404            {object}  resp.ErrorResponse  "Публикация не найдена"
// @Failure      500            {object}  resp.ErrorResponse  "Внутренняя ошибка сервера"
// @Router       /publications/{publicationID}/form [post]
func (h *Handler) HandleFormPublication(c *gin.Context) {
	userID, err := getUserIDFromContext(c)
	if err != nil {
		resp.WriteError(c.Writer, http.StatusUnauthorized, resp.SingleError(err.Error()), nil)
		return
	}

	id, err := strconv.Atoi(c.Param("publicationID"))
	if err != nil {
		resp.WriteError(c.Writer, http.StatusBadRequest, resp.ErrorDetailList("publicationID", err.Error()), nil)
		return
	}

	publication, err := h.repo.GetPublicationByID(userID, id)
	if err != nil {
		if errors.Is(err, models.ErrPublicationNotFound) {
			resp.WriteError(c.Writer, http.StatusNotFound, resp.SingleError(err.Error()), nil)
		} else {
			resp.WriteError(c.Writer, http.StatusInternalServerError, resp.SingleError(err.Error()), nil)
		}
		return
	}

	if publication.Status != models.DraftPublicationStatus && publication.UserID != userID {
		resp.WriteError(c.Writer, http.StatusConflict, resp.ErrorDetailList("status", "status is not draft"), nil)
		return
	}

	publication, err = h.repo.UpdatePublicationStatus(userID, id, models.WorkPublicationStatus)
	if err != nil {
		resp.WriteError(c.Writer, http.StatusInternalServerError, resp.SingleError(err.Error()), nil)
		return
	}

	resp.WriteJSON(c.Writer, http.StatusOK, publication)
}

// HandleFinalizedPublication godoc
// @Summary      Завершить публикацию
// @Description  Переводит публикацию в финальный статус.
// @Tags         Publications
// @Accept       json
// @Produce      json
// @Param        publicationID  path    int  true  "Идентификатор публикации"
// @Param        status         query   string  true  "Статус публикации"
// @Security BearerAuth
// @Success      200            {object}  models.Publication   "Завершенная публикация"
// @Failure      400            {object}  resp.ErrorResponse  "Неверный формат данных"
// @Failure      401            {object}  resp.ErrorResponse  "Неверные учетные данные"
// @Failure 	 403 	  {object}  resp.ErrorResponse  "Недостаточно прав"
// @Failure      404            {object}  resp.ErrorResponse  "Публикация не найдена"
// @Failure      409            {object}  resp.ErrorResponse  "Неверный статус для завершения публикации"
// @Failure      500            {object}  resp.ErrorResponse  "Внутренняя ошибка сервера"
// @Router       /publications/{publicationID}/finalize [post]
func (h *Handler) HandleFinalizedPublication(c *gin.Context) {
	userID, err := getUserIDFromContext(c)
	if err != nil {
		resp.WriteError(c.Writer, http.StatusUnauthorized, resp.SingleError(err.Error()), nil)
		return
	}

	id, err := strconv.Atoi(c.Param("publicationID"))
	if err != nil {
		resp.WriteError(c.Writer, http.StatusBadRequest, resp.ErrorDetailList("publicationID", err.Error()), nil)
		return
	}

	status := models.PublicationStatus(c.Query("status"))

	publication, err := h.repo.GetPublicationByID(userID, id)
	if err != nil {
		resp.WriteError(c.Writer, http.StatusInternalServerError, resp.SingleError(err.Error()), nil)
		return
	}

	if publication.Status != models.WorkPublicationStatus {
		resp.WriteError(c.Writer, http.StatusConflict, resp.ErrorDetailList("status", "status is not work"), nil)
		return
	}

	publication, err = h.repo.UpdatePublicationStatus(userID, id, status)
	if err != nil {
		resp.WriteError(c.Writer, http.StatusInternalServerError, resp.SingleError(err.Error()), nil)
		return
	}

	resp.WriteJSON(c.Writer, http.StatusOK, publication)
}

// HandleDeletePublication godoc
// @Summary      Удалить публикацию по ID
// @Description  Удаляет публикацию по уникальному идентификатору.
// @Tags         Publications
// @Accept       json
// @Produce      json
// @Param        publicationID  path    int  true  "Идентификатор публикации"
// @Security BearerAuth
// @Success      204            "Публикация успешно удалена"
// @Failure      400            {object}  resp.ErrorResponse  "Неверный формат ID публикации"
// @Failure      401            {object}  resp.ErrorResponse  "Неверные учетные данные"
// @Failure      404            {object}  resp.ErrorResponse  "Публикация не найдена"
// @Failure      500            {object}  resp.ErrorResponse  "Внутренняя ошибка сервера"
// @Router       /publications/{publicationID} [delete]
func (h *Handler) HandleDeletePublication(c *gin.Context) {
	userID, err := getUserIDFromContext(c)
	if err != nil {
		resp.WriteError(c.Writer, http.StatusUnauthorized, resp.SingleError(err.Error()), nil)
		return
	}

	id, err := strconv.Atoi(c.Param("publicationID"))
	if err != nil {
		resp.WriteError(c.Writer, http.StatusBadRequest, resp.ErrorDetailList("publicationID", err.Error()), nil)
		return
	}

	err = h.repo.DeletePublication(userID, id)
	if err != nil {
		if errors.Is(err, models.ErrPublicationNotFound) {
			resp.WriteError(c.Writer, http.StatusNotFound, resp.SingleError(err.Error()), nil)
		} else {
			resp.WriteError(c.Writer, http.StatusInternalServerError, resp.SingleError(err.Error()), nil)
		}
		return
	}

	resp.WriteJSON(c.Writer, http.StatusNoContent, nil)
}
