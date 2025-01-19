package postgresrepo

import (
	"backend/internal/models"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

func (r *PostgresRepository) RemoveEventFromPublication(userID int, eventID int) error {
	var event models.HistoricalEvent
	if err := r.db.Where("id = ? AND status = ?", eventID, models.ActiveEventStatus).First(&event).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return models.ErrPublicationNotFound
		}
		return errors.Wrap(err, "failed to fetch publication")
	}

	var publication models.Publication
	if err := r.db.Where("user_id = ? AND status = ?", userID, models.DraftPublicationStatus).First(&publication).Error; err != nil {
		return models.ErrPublicationDraftNotFound
	}

	var publicationEvent models.PublicationsEvents
	if err := r.db.Where("publication_id = ? AND event_id = ?", publication.ID, eventID).First(&publicationEvent).Error; err != nil {
		return errors.Wrap(err, "event not found in publication")
	}

	if err := r.db.Delete(&publicationEvent).Error; err != nil {
		return errors.Wrap(err, "failed to remove event from publication")
	}

	return nil
}

func (r *PostgresRepository) UpdateEventPriority(userID int, publicationID, eventID int, priority int) error {
	var publication models.Publication
	if err := r.db.Where("id = ? AND user_id = ?", publicationID, userID).First(&publication).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return models.ErrPublicationNotFound
		}
		return errors.Wrap(err, "failed to fetch publication")
	}

	var publicationEvent models.PublicationsEvents
	if err := r.db.Where("publication_id = ? AND event_id = ?", publicationID, eventID).First(&publicationEvent).Error; err != nil {
		return errors.Wrap(err, "event not found in publication")
	}

	publicationEvent.Priority = priority
	if err := r.db.Save(&publicationEvent).Error; err != nil {
		return errors.Wrap(err, "failed to update event priority")
	}

	return nil
}

func calculateTrustScore(publication *models.Publication) float64 {
	baseScore := 0.0

	// Установить низкий рейтинг, если публикация отклонена
	if publication.Status == models.RejectedPublicationStatus {
		return 10 // Минимальный базовый рейтинг для отклонённых публикаций
	}

	// Базовая логика: если публикация завершена, добавляем базовый балл
	if publication.Status == models.CompletedPublicationStatus {
		baseScore += 50
	}

	// Учитываем время между созданием и завершением
	if publication.FormationDate != nil && publication.CompletionDate != nil {
		duration := publication.CompletionDate.Sub(*publication.FormationDate).Hours() / 24 // дни
		if duration < 7 {
			baseScore += 20 // Быстрая модерация
		} else if duration <= 30 {
			baseScore += 10 // Средний срок
		} else {
			baseScore -= 10 // Затянувшаяся модерация
		}
	}

	// Учитываем наличие описания
	if publication.Description != "" && publication.Description != "Это описание черновой публикации" {
		baseScore += 10 // Полнота информации
	}

	// Учитываем длину названия
	if len(publication.Title) > 10 {
		baseScore += 5 // Хорошее название
	} else {
		baseScore -= 5 // Слишком короткое название
	}

	// Минимальный и максимальный лимиты
	if baseScore > 100 {
		baseScore = 100
	} else if baseScore < 0 {
		baseScore = 0
	}

	return baseScore
}
