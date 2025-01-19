package handlers

import (
	"backend/internal/models"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"net/http"
	"strconv"
)

func (h *Handler) HandleGetEvents(c *gin.Context) {
	userID, err := getUserIDFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

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

	draftPublication, _ := h.repo.GetDraftPublication(userID)
	var draftPublicationID, draftPublicationEventsCount int
	if draftPublication != nil {
		draftPublicationID = draftPublication.ID
		draftPublicationEventsCount = len(draftPublication.Events)
	}

	c.JSON(http.StatusOK, models.GetEventsDTO{
		HistoricalEvents: events,
		PublicationID:    draftPublicationID,
		EventsCount:      draftPublicationEventsCount,
	})
}

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
