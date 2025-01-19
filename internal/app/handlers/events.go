package handlers

import (
	"backend/internal/models"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"net/http"
	"strconv"
)

// HandleGetEvents godoc
// @Summary      Получить список событий
// @Description  Возвращает список событий с возможностью фильтрации по названию. Также возвращает черновик публикации пользователя, если он существует.
// @Tags         Events
// @Accept       json
// @Produce      json
// @Param        title  query    string  false  "Название события"
// @Success      200    {object}  models.GetEventsDTO
// @Failure      404    {object}  gin.H  "События не найдены"
// @Failure      500    {object}  gin.H  "Внутренняя ошибка сервера"
// @Router       /events [get]
func (h *Handler) HandleGetEvents(c *gin.Context) {
	title := c.DefaultQuery("title", "")
	events, err := h.repo.GetEvents(title)
	if err != nil {
		if errors.Is(err, models.ErrEventNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	var draftPublicationID, draftPublicationEventsCount int

	userID, err := getUserIDFromContext(c)
	if err == nil {
		draftPublication, _ := h.repo.GetDraftPublication(userID)
		if draftPublication != nil {
			draftPublicationID = draftPublication.ID
			draftPublicationEventsCount = len(draftPublication.Events)
		}
	}

	c.JSON(http.StatusOK, models.GetEventsDTO{
		HistoricalEvents: events,
		PublicationID:    draftPublicationID,
		EventsCount:      draftPublicationEventsCount,
	})
}

// HandleGetEventByID godoc
// @Summary      Получить событие по ID
// @Description  Возвращает событие по уникальному идентификатору.
// @Tags         Events
// @Accept       json
// @Produce      json
// @Param        eventID  path    int  true  "Идентификатор события"
// @Success      200      {object}  models.HistoricalEvent
// @Failure      400      {object}  gin.H  "Неверный формат ID события"
// @Failure      404      {object}  gin.H  "Событие не найдено"
// @Failure      500      {object}  gin.H  "Внутренняя ошибка сервера"
// @Router       /events/{eventID} [get]
func (h *Handler) HandleGetEventByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("eventID"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid event ID"})
		return
	}

	event, err := h.repo.GetEventByID(id)
	if err != nil {
		if errors.Is(err, models.ErrEventNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, event)
}

// HandleAddEvent godoc
// @Summary      Добавить новое событие
// @Description  Добавляет новое событие в систему.
// @Tags         Events
// @Accept       json
// @Produce      json
// @Param        event  body    models.CreateEventDTO  true  "Информация о событии"
// @Security BearerAuth
// @Success      201     {object}  models.HistoricalEvent
// @Failure      400     {object}  gin.H  "Неверный формат данных"
// @Failure      401      {object}  gin.H  "Неверные учетные данные"
// @Failure 	 403 	  {object}  gin.H  "Недостаточно прав"
// @Failure      500     {object}  gin.H  "Внутренняя ошибка сервера"
// @Router       /events [post]
func (h *Handler) HandleAddEvent(c *gin.Context) {
	var event models.CreateEventDTO
	if err := c.ShouldBindJSON(&event); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	newEvent, err := h.repo.CreateEvent(&event)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, newEvent)
}

// HandleUpdateEvent godoc
// @Summary      Обновить событие по ID
// @Description  Обновляет информацию о событии по его уникальному идентификатору.
// @Tags         Events
// @Accept       json
// @Produce      json
// @Param        eventID  path    int                        true  "Идентификатор события"
// @Param        event    body    models.UpdateEventDTO      true  "Обновленные данные события"
// @Security BearerAuth
// @Success      200      {object}  models.HistoricalEvent
// @Failure      400      {object}  gin.H  "Неверный формат данных"
// @Failure      401      {object}  gin.H  "Неверные учетные данные"
// @Failure 	 403 	  {object}  gin.H  "Недостаточно прав"
// @Failure      404      {object}  gin.H  "Событие не найдено"
// @Failure      409      {object}  gin.H  "Событие уже удалено"
// @Failure      500      {object}  gin.H  "Внутренняя ошибка сервера"
// @Router       /events/{eventID} [put]
func (h *Handler) HandleUpdateEvent(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("eventID"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid event ID"})
		return
	}

	var event models.UpdateEventDTO
	if err = c.ShouldBindJSON(&event); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	updatedEvent, err := h.repo.UpdateEvent(id, &event)
	if err != nil {
		if errors.Is(err, models.ErrEventNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		} else if errors.Is(err, models.ErrEventAlreadyDeleted) {
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, updatedEvent)
}

// HandleDeleteEvent godoc
// @Summary      Удалить событие по ID
// @Description  Удаляет событие по его уникальному идентификатору.
// @Tags         Events
// @Accept       json
// @Produce      json
// @Param        eventID  path    int  true  "Идентификатор события"
// @Security BearerAuth
// @Success      204      {object}  nil
// @Failure      400      {object}  gin.H  "Неверный формат ID события"
// @Failure      401      {object}  gin.H  "Неверные учетные данные"
// @Failure 	 403 	  {object}  gin.H  "Недостаточно прав"
// @Failure      404      {object}  gin.H  "Событие не найдено"
// @Failure      409      {object}  gin.H  "Событие уже удалено"
// @Failure      500      {object}  gin.H  "Внутренняя ошибка сервера"
// @Router       /events/{eventID} [delete]
func (h *Handler) HandleDeleteEvent(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("eventID"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid event ID"})
		return
	}

	err = h.repo.DeleteEvent(id)
	if err != nil {
		if errors.Is(err, models.ErrEventNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		} else if errors.Is(err, models.ErrEventAlreadyDeleted) {
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	c.JSON(http.StatusNoContent, nil)
}

// HandleAddEventToPublication godoc
// @Summary      Добавить событие в публикацию
// @Description  Добавляет событие в черновик публикации пользователя.
// @Tags         Events
// @Accept       json
// @Produce      json
// @Param        eventID  path    int  true  "Идентификатор события"
// @Security BearerAuth
// @Success      204      {object}  nil
// @Failure      400      {object}  gin.H  "Неверный формат ID события"
// @Failure      401      {object}  gin.H  "Неверные учетные данные"
// @Failure      404      {object}  gin.H  "Событие или черновик публикации не найдено"
// @Failure      500      {object}  gin.H  "Внутренняя ошибка сервера"
// @Router       /events/{eventID}/publications [post]
func (h *Handler) HandleAddEventToPublication(c *gin.Context) {
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

	err = h.repo.AddEventToPublication(userID, eventID)
	if err != nil {
		if errors.Is(err, models.ErrEventNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		} else if errors.Is(err, models.ErrPublicationDraftNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	c.JSON(http.StatusNoContent, nil)
}

// HandleAddEventImage godoc
// @Summary      Добавить изображение к событию
// @Description  Добавляет изображение к событию по его ID.
// @Tags         Events
// @Accept       multipart/form-data
// @Produce      json
// @Param        eventID  path    int         true  "Идентификатор события"
// @Param        image    formData  file    true  "Изображение события"
// @Security BearerAuth
// @Success      201      {object}  nil
// @Failure      400      {object}  gin.H  "Неверный формат ID события или файла"
// @Failure      401      {object}  gin.H  "Неверные учетные данные"
// @Failure 	 403 	  {object}  gin.H  "Недостаточно прав"
// @Failure      404      {object}  gin.H  "Событие не найдено"
// @Failure      500      {object}  gin.H  "Внутренняя ошибка сервера"
// @Router       /events/{eventID}/image [post]
func (h *Handler) HandleAddEventImage(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("eventID"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid event ID"})
		return
	}

	file, err := c.FormFile("image")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid file"})
		return
	}

	err = h.repo.AddImageToEvent(id, file)
	if err != nil {
		if errors.Is(err, models.ErrEventNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	c.JSON(http.StatusCreated, nil)
}
