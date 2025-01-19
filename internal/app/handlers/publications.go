package handlers

import (
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
// @Failure      400      {object}  gin.H  "Неверный формат параметров"
// @Failure      401      {object}  gin.H  "Неверные учетные данные"
// @Failure      404      {object}  gin.H  "Публикации не найдены"
// @Failure      500      {object}  gin.H  "Внутренняя ошибка сервера"
// @Router       /publications [get]
func (h *Handler) HandleGetPublications(c *gin.Context) {
	userID, err := getUserIDFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	status := models.PublicationStatus(c.Query("status"))

	startDateStr := c.Query("startDate")
	endDateStr := c.Query("endDate")

	var startDate, endDate *time.Time
	if startDateStr != "" {
		parsedStartDate, err := time.Parse("2006-01-02", startDateStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid start date format"})
			return
		}
		startDate = &parsedStartDate
	}

	if endDateStr != "" {
		parsedEndDate, err := time.Parse("2006-01-02", endDateStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid end date format"})
			return
		}
		endDate = &parsedEndDate
	}

	publications, err := h.repo.GetPublications(userID, status, startDate, endDate)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if len(publications) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, publications)
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
// @Failure      400            {object}  gin.H  "Неверный формат ID публикации"
// @Failure      401            {object}  gin.H  "Неверные учетные данные"
// @Failure      404            {object}  gin.H  "Публикация не найдена"
// @Failure      500            {object}  gin.H  "Внутренняя ошибка сервера"
// @Router       /publications/{publicationID} [get]
func (h *Handler) HandleGetPublicationByID(c *gin.Context) {
	userID, err := getUserIDFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	id, err := strconv.Atoi(c.Param("publicationID"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid publication ID"})
		return
	}

	publication, err := h.repo.GetPublicationByID(userID, id)
	if err != nil {
		if errors.Is(err, models.ErrPublicationNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "publication not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, publication)
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
// @Failure      400            {object}  gin.H  "Неверный формат данных"
// @Failure      401            {object}  gin.H  "Неверные учетные данные"
// @Failure      404            {object}  gin.H  "Публикация не найдена"
// @Failure      500            {object}  gin.H  "Внутренняя ошибка сервера"
// @Router       /publications/{publicationID} [put]
func (h *Handler) HandleUpdatePublication(c *gin.Context) {
	userID, err := getUserIDFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	id, err := strconv.Atoi(c.Param("publicationID"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid publication ID"})
		return
	}

	var publication models.UpdatePublicationDTO
	if err := c.ShouldBindJSON(&publication); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	updatedPublication, err := h.repo.UpdatePublication(userID, id, &publication)
	if err != nil {
		if errors.Is(err, models.ErrPublicationNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, updatedPublication)
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
// @Failure      400            {object}  gin.H  "Неверный формат ID публикации"
// @Failure      401            {object}  gin.H  "Неверные учетные данные"
// @Failure      403            {object}  gin.H  "Недостаточно прав"
// @Failure      404            {object}  gin.H  "Публикация не найдена"
// @Failure      500            {object}  gin.H  "Внутренняя ошибка сервера"
// @Router       /publications/{publicationID}/form [post]
func (h *Handler) HandleFormPublication(c *gin.Context) {
	userID, err := getUserIDFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	id, err := strconv.Atoi(c.Param("publicationID"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid publication ID"})
		return
	}

	publication, err := h.repo.GetPublicationByID(userID, id)
	if err != nil {
		if errors.Is(err, models.ErrPublicationNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	if publication.Status != models.DraftPublicationStatus && publication.UserID != userID {
		c.JSON(http.StatusForbidden, gin.H{"error": "you are not authorized to edit this publication"})
		return
	}

	publication, err = h.repo.UpdatePublicationStatus(userID, id, models.WorkPublicationStatus)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, publication)
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
// @Failure      400            {object}  gin.H  "Неверный формат данных"
// @Failure      401            {object}  gin.H  "Неверные учетные данные"
// @Failure 	 403 	  {object}  gin.H  "Недостаточно прав"
// @Failure      404            {object}  gin.H  "Публикация не найдена"
// @Failure      409            {object}  gin.H  "Неверный статус для завершения публикации"
// @Failure      500            {object}  gin.H  "Внутренняя ошибка сервера"
// @Router       /publications/{publicationID}/finalize [post]
func (h *Handler) HandleFinalizedPublication(c *gin.Context) {
	userID, err := getUserIDFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	id, err := strconv.Atoi(c.Param("publicationID"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid publication ID"})
		return
	}

	status := models.PublicationStatus(c.Query("status"))

	publication, err := h.repo.GetPublicationByID(userID, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if publication.Status != models.WorkPublicationStatus {
		c.JSON(http.StatusConflict, gin.H{"error": "invalid publication status for finalization"})
		return
	}

	publication, err = h.repo.UpdatePublicationStatus(userID, id, status)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, publication)
}

// HandleDeletePublication godoc
// @Summary      Удалить публикацию по ID
// @Description  Удаляет публикацию по уникальному идентификатору.
// @Tags         Publications
// @Accept       json
// @Produce      json
// @Param        publicationID  path    int  true  "Идентификатор публикации"
// @Security BearerAuth
// @Success      204            {object}  nil   "Публикация успешно удалена"
// @Failure      400            {object}  gin.H  "Неверный формат ID публикации"
// @Failure      401            {object}  gin.H  "Неверные учетные данные"
// @Failure      404            {object}  gin.H  "Публикация не найдена"
// @Failure      500            {object}  gin.H  "Внутренняя ошибка сервера"
// @Router       /publications/{publicationID} [delete]
func (h *Handler) HandleDeletePublication(c *gin.Context) {
	userID, err := getUserIDFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	id, err := strconv.Atoi(c.Param("publicationID"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid publication ID"})
		return
	}

	err = h.repo.DeletePublication(userID, id)
	if err != nil {
		if errors.Is(err, models.ErrPublicationNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	c.JSON(http.StatusNoContent, nil)
}
