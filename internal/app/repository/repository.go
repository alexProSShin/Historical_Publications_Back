package repository

import (
	"backend/internal/models"
	"github.com/pkg/errors"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"strings"
)

type Repository interface {
	GetEventByID(id int) (*models.HistoricalEvent, error)
	GetEvents(title string) ([]models.HistoricalEvent, error)
	DeleteEvent(id int) error
	AddEventToPublication(eventID int) error
	RemoveEventFromPublication(eventID int) error
	GetDraftPublication() (*models.GetPublicationDTO, error)
	GetPublicationByID(id int) (*models.GetPublicationDTO, error)
	UpdatePublicationStatusByID(id int, status models.PublicationStatus) error
}

type PostgresRepository struct {
	db *gorm.DB
}

func NewPostgresRepository(dsn string) (*PostgresRepository, error) {
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	return &PostgresRepository{
		db: db,
	}, nil
}

func (r *PostgresRepository) GetEventByID(id int) (*models.HistoricalEvent, error) {
	var event models.HistoricalEvent
	err := r.db.First(&event, "id = ? AND status = ?", id, models.ActiveEventStatus).Error
	if err != nil {
		return nil, errors.Wrapf(err, "failed to get event by id %d", id)
	}
	return &event, nil
}

func (r *PostgresRepository) GetEvents(title string) ([]models.HistoricalEvent, error) {
	var events []models.HistoricalEvent
	title = strings.ToLower(title + "%")
	if err := r.db.Find(&events, "status = ? AND LOWER(title) LIKE ?", models.ActiveEventStatus, title).Error; err != nil {
		return nil, errors.Wrapf(err, "failed to get events by title %s", title)
	}

	return events, nil
}

func (r *PostgresRepository) DeleteEvent(id int) error {
	return r.db.Exec("UPDATE ships SET status = ? WHERE id = ?", models.ActiveEventStatus, id).Error
}

func (r *PostgresRepository) AddEventToPublication(eventID int) error {
	var event models.HistoricalEvent

	if err := r.db.
		Where("id = ? AND status = ?", eventID, models.ActiveEventStatus).
		First(&event).Error; err != nil {
		return errors.Wrap(err, models.ErrEventNotFound)
	}

	var publication models.Publication

	if err := r.db.
		Where("status = ?", models.DraftPublicationStatus).
		Last(&publication).Error; err != nil {
		if err = r.db.
			Create(&publication).Error; err != nil {
			return errors.Wrap(err, "ошибка создания заявки со статусом черновик")
		}
	}

	var existingPublicationEvent models.PublicationsEvents
	if err := r.db.
		Where("event_id = ? AND publication_id = ?", event.ID, publication.ID).
		First(&existingPublicationEvent).Error; err == nil {
		return nil
	}

	publicationEvent := models.PublicationsEvents{
		EventID:       event.ID,
		PublicationID: publication.ID,
	}

	if err := r.db.
		Create(&publicationEvent).Error; err != nil {
		return errors.Wrap(err, "ошибка при создании связи между заявкой и событием")
	}

	return nil
}

func (r *PostgresRepository) RemoveEventFromPublication(eventID int) error {
	var publicationEvent models.PublicationsEvents

	if err := r.db.Joins("JOIN publications ON publications_events.request_id = publications.id").
		Where("publications_events.event_id = ? AND publications.status = ?", eventID, models.DraftPublicationStatus).
		First(&publicationEvent).Error; err != nil {
		return errors.Wrap(err, "судно не принадлежит пользователю или находится не в статусе черновик")
	}

	if err := r.db.
		Delete(&publicationEvent).Error; err != nil {
		return errors.Wrap(err, "ошибка удаления связи между заявкой и судном")
	}

	return nil
}

func (r *PostgresRepository) GetDraftPublication() (*models.GetPublicationDTO, error) {
	var publication models.Publication
	err := r.db.Where("status = ?", models.DraftPublicationStatus).Last(&publication).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, errors.Wrap(err, "не удалось получить публикацию черновика")
	}

	var events []models.HistoricalEvent
	err = r.db.Joins("JOIN publications_events ON publications_events.event_id = historical_events.id").
		Where("publications_events.publication_id = ?", publication.ID).
		Find(&events).Error
	if err != nil {
		return nil, errors.Wrapf(err, "не удалось получить события для публикации с id %d", publication.ID)
	}

	publicationDTO := &models.GetPublicationDTO{
		Publication: publication,
		Events:      events,
	}

	return publicationDTO, nil
}

func (r *PostgresRepository) GetPublicationByID(id int) (*models.GetPublicationDTO, error) {
	var publication models.Publication
	err := r.db.First(&publication, "id = ?", id).Error
	if err != nil {
		return nil, errors.Wrapf(err, "не удалось получить публикацию с id %d", id)
	}

	var events []models.HistoricalEvent
	err = r.db.Joins("JOIN publications_events ON publications_events.event_id = historical_events.id").
		Where("publications_events.publication_id = ?", id).
		Find(&events).Error
	if err != nil {
		return nil, errors.Wrapf(err, "не удалось получить события для публикации с id %d", id)
	}

	publicationDTO := &models.GetPublicationDTO{
		Publication: publication,
		Events:      events,
	}

	return publicationDTO, nil
}

func (r *PostgresRepository) UpdatePublicationStatusByID(id int, status models.PublicationStatus) error {
	var publication models.Publication
	err := r.db.First(&publication, "id = ?", id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.Wrapf(err, "публикация с id %d не найдена", id)
		}
		return errors.Wrapf(err, "не удалось найти публикацию с id %d", id)
	}

	publication.Status = status
	err = r.db.Save(&publication).Error
	if err != nil {
		return errors.Wrapf(err, "не удалось обновить статус публикации с id %d", id)
	}

	return nil
}
