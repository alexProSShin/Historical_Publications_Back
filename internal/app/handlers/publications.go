package handlers

import (
	"backend/internal/models"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"net/http"
	"strconv"
	"time"
)

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
