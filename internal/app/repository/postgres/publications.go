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
	var events []models.GetEvents
	if err := r.db.Table("historical_events").
		Select("historical_events.*, publications_events.priority AS priority").
		Joins("JOIN publications_events ON publications_events.event_id = historical_events.id").
		Where("publications_events.publication_id = ?", publication.ID).
		Order("publications_events.priority DESC").
		Scan(&events).Error; err != nil {
		return nil, errors.Wrap(err, "failed to fetch events with priorities for publication")
	}
	return &models.GetPublicationDTO{
		Publication: publication,
		Events:      events,
	}, nil
}

func (r *PostgresRepository) GetPublications(user *models.User, status models.PublicationStatus, startDate, endDate *time.Time) ([]models.GetPublications, error) {
	var publications []models.GetPublications

	query := r.db.Table("publications").
		Select("publications.*, users.name AS user_name").
		Joins("LEFT JOIN users ON users.id = publications.user_id")

	if user.Role == models.RoleUser {
		query = query.Where("publications.user_id = ? AND publications.status != ?", user.ID, models.DeletedPublicationStatus)
	} else if user.Role == models.RoleModerator {
		query = query.Where("publications.status != ? AND publications.status != ?", models.DeletedPublicationStatus, models.DraftPublicationStatus)
	}

	if status != "" {
		query = query.Where("publications.status = ?", status)
	}

	if startDate != nil {
		query = query.Where("publications.formation_date >= ?", *startDate)
	}
	if endDate != nil {
		query = query.Where("publications.formation_date <= ?", *endDate)
	}

	if err := query.Find(&publications).Error; err != nil {
		return nil, errors.Wrap(err, "failed to fetch publications with user names")
	}

	return publications, nil
}

func (r *PostgresRepository) GetPublicationByID(user *models.User, id int) (*models.GetPublicationDTO, error) {
	var publication models.Publication
	query := r.db.Where("id = ?", id)

	if user.Role == models.RoleUser {
		query = query.Where("user_id = ?", user.ID)
	}

	if err := query.First(&publication).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, models.ErrPublicationNotFound
		}
		return nil, errors.Wrap(err, "failed to fetch publication")
	}

	var events []models.GetEvents
	if err := r.db.Table("historical_events").
		Select("historical_events.*, publications_events.priority AS priority").
		Joins("JOIN publications_events ON publications_events.event_id = historical_events.id").
		Where("publications_events.publication_id = ?", id).
		Order("publications_events.priority DESC").
		Scan(&events).Error; err != nil {
		return nil, errors.Wrap(err, "failed to fetch events with priorities for publication")
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

	user, err := r.GetUserByID(userID)
	if err != nil {
		return nil, errors.Wrap(err, "failed to fetch user")
	}

	if err := r.db.Model(&existingPublication).Updates(publication).Error; err != nil {
		return nil, errors.Wrap(err, "failed to update publication")
	}

	return r.GetPublicationByID(user, id)
}

func (r *PostgresRepository) UpdatePublicationStatus(user *models.User, id int, status models.PublicationStatus) (*models.GetPublicationDTO, error) {
	var publication models.Publication
	query := r.db.Where("id = ?", id)

	if user.Role == models.RoleUser {
		query = query.Where("user_id = ?", user.ID)
	}

	if err := query.First(&publication).Error; err != nil {
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
		publication.ModeratorID = &user.ID
		trustScore := calculateTrustScore(&publication)
		publication.TrustScore = &trustScore
	}

	publication.Status = status
	if err := r.db.Save(&publication).Error; err != nil {
		return nil, errors.Wrap(err, "failed to update publication status")
	}

	return r.GetPublicationByID(user, id)
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
