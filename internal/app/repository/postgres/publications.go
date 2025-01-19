package postgresrepo

import (
	"backend/internal/models"
	"github.com/pkg/errors"
	"gorm.io/gorm"
	"time"
)

func (r *PostgresRepository) GetDraftPublication(userID int) (*models.GetPublicationDTO, error) {
	var publication models.Publication
	err := r.db.Where("status = ? AND user_id = ?", models.DraftPublicationStatus, userID).Last(&publication).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, models.ErrPublicationDraftNotFound
		}
		return nil, errors.Wrap(err, "failed to fetch draft publication")
	}
	var events []models.HistoricalEvent
	err = r.db.Joins("JOIN publications_events ON publications_events.event_id = historical_events.id").
		Where("publications_events.publication_id = ?", publication.ID).
		Find(&events).Error
	if err != nil {
		return nil, errors.Wrap(err, "failed to fetch events for draft publication")
	}
	return &models.GetPublicationDTO{
		Publication: publication,
		Events:      events,
	}, nil
}

func (r *PostgresRepository) GetPublications(userID int, status models.PublicationStatus, startDate, endDate *time.Time) ([]models.Publication, error) {
	var publications []models.Publication
	query := r.db.Where("user_id = ? AND status != ?", userID, models.DeletedPublicationStatus)

	if status != "" {
		query = query.Where("status = ?", status)
	}
	if startDate != nil {
		query = query.Where("formation_date >= ?", *startDate)
	}
	if endDate != nil {
		query = query.Where("formation_date <= ?", *endDate)
	}

	if err := query.Find(&publications).Error; err != nil {
		return nil, errors.Wrap(err, "failed to fetch publications")
	}
	return publications, nil
}

func (r *PostgresRepository) GetPublicationByID(userID int, id int) (*models.GetPublicationDTO, error) {
	var publication models.Publication
	if err := r.db.Where("id = ? AND user_id = ?", id, userID).First(&publication).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, models.ErrPublicationNotFound
		}
		return nil, errors.Wrap(err, "failed to fetch publication")
	}

	var events []models.HistoricalEvent
	if err := r.db.Joins("JOIN publications_events ON publications_events.event_id = historical_events.id").
		Where("publications_events.publication_id = ?", id).
		Order("publications_events.priority DESC").
		Find(&events).Error; err != nil {
		return nil, errors.Wrap(err, "failed to fetch events for publication")
	}

	return &models.GetPublicationDTO{
		Publication: publication,
		Events:      events,
	}, nil
}

func (r *PostgresRepository) UpdatePublication(userID int, id int, publication *models.UpdatePublicationDTO) (*models.GetPublicationDTO, error) {
	var existingPublication models.Publication
	if err := r.db.Where("id = ? AND user_id = ?", id, userID).First(&existingPublication).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, models.ErrPublicationNotFound
		}
		return nil, errors.Wrap(err, "failed to fetch publication")
	}

	if err := r.db.Model(&existingPublication).Updates(publication).Error; err != nil {
		return nil, errors.Wrap(err, "failed to update publication")
	}

	return r.GetPublicationByID(userID, id)
}

func (r *PostgresRepository) UpdatePublicationStatus(userID, id int, status models.PublicationStatus) (*models.GetPublicationDTO, error) {
	var publication models.Publication
	if err := r.db.Where("id = ? AND user_id = ?", id, userID).First(&publication).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, models.ErrPublicationNotFound
		}
		return nil, errors.Wrap(err, "failed to fetch publication")
	}

	if publication.Status == models.DeletedPublicationStatus {
		return nil, models.ErrPublicationAlreadyDone
	}

	currentTime := time.Now()
	if status == models.WorkPublicationStatus {
		publication.FormationDate = &currentTime
	} else {
		publication.CompletionDate = &currentTime
		publication.ModeratorID = &userID
		trustScore := calculateTrustScore(&publication)
		publication.TrustScore = &trustScore
	}

	publication.Status = status
	if err := r.db.Save(&publication).Error; err != nil {
		return nil, errors.Wrap(err, "failed to update publication status")
	}

	return r.GetPublicationByID(userID, id)
}

func (r *PostgresRepository) DeletePublication(userID int, id int) error {
	var publication models.Publication
	if err := r.db.Where("id = ? AND user_id = ?", id, userID).First(&publication).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return models.ErrPublicationNotFound
		}
		return errors.Wrap(err, "failed to fetch publication")
	}

	if publication.Status == models.DeletedPublicationStatus {
		return models.ErrPublicationAlreadyDone
	}

	if err := r.db.Model(&publication).Update("status", models.DeletedPublicationStatus).Error; err != nil {
		return errors.Wrap(err, "failed to delete publication")
	}

	return nil
}
