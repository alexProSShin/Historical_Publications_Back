package handlers

import (
	"backend/internal/models"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"net/http"
	"strconv"
)

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
