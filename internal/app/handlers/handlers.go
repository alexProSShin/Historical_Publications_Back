package handlers

import (
	"backend/internal/app/repository"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type Handler struct {
	repo repository.Repository
}

func New(repo repository.Repository) *Handler {
	return &Handler{
		repo: repo,
	}
}

func (h *Handler) HandleGetEvents(c *gin.Context) {
	title := c.DefaultQuery("title", "")
	events, err := h.repo.GetEvents(title)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	draftPublication, _ := h.repo.GetDraftPublication()
	var draftPublicationID, draftPublicationEventsCount int
	if draftPublication != nil {
		draftPublicationID = draftPublication.ID
		draftPublicationEventsCount = len(draftPublication.Events)
	}

	c.HTML(http.StatusOK, "main_page.tmpl", gin.H{
		"events":        events,
		"searchValue":   title,
		"publicationId": draftPublicationID,
		"eventsCount":   draftPublicationEventsCount,
	})
}

func (h *Handler) HandleGetEventByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid event ID"})
		return
	}

	event, err := h.repo.GetEventByID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	draftPublication, _ := h.repo.GetDraftPublication()
	var draftPublicationID, draftPublicationEventsCount int
	if draftPublication != nil {
		draftPublicationID = draftPublication.ID
		draftPublicationEventsCount = len(draftPublication.Events)
	}

	c.HTML(http.StatusOK, "event_page.tmpl", gin.H{
		"event":         event,
		"publicationId": draftPublicationID,
		"eventsCount":   draftPublicationEventsCount,
	})
}

func (h *Handler) HandleAddEventToPublication(c *gin.Context) {
	eventID, err := strconv.Atoi(c.Param("eventID"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid event ID"})
		return
	}

	err = h.repo.AddEventToPublication(eventID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Redirect(http.StatusFound, "/")
}

func (h *Handler) HandleGetPublicationByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid publication ID"})
		return
	}

	publication, err := h.repo.GetPublicationByID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.HTML(http.StatusOK, "publication_page.tmpl", gin.H{
		"publication": publication,
		"eventsCount": len(publication.Events),
	})
}
