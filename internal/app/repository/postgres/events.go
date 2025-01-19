package postgresrepo

import (
	"backend/internal/models"
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/minio/minio-go/v7"
	"github.com/pkg/errors"
	"gorm.io/gorm"
	"mime/multipart"
	"strings"
)

func (r *PostgresRepository) GetEvents(title string) ([]models.HistoricalEvent, error) {
	var events []models.HistoricalEvent
	title = strings.ToLower(title) + "%"
	if err := r.db.Where("status = ? AND LOWER(title) LIKE ?", models.ActiveEventStatus, title).Find(&events).Error; err != nil {
		return nil, errors.Wrap(err, "failed to fetch events")
	}
	return events, nil
}

func (r *PostgresRepository) GetEventByID(id int) (*models.HistoricalEvent, error) {
	var event models.HistoricalEvent
	if err := r.db.Where("id = ? AND status = ?", id, models.ActiveEventStatus).First(&event).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, models.ErrEventNotFound
		}
		return nil, errors.Wrap(err, "failed to fetch event")
	}
	return &event, nil
}

func (r *PostgresRepository) CreateEvent(event *models.CreateEventDTO) (*models.HistoricalEvent, error) {
	historicalEvent := &models.HistoricalEvent{
		Title:       event.Title,
		EventType:   event.EventType,
		Description: event.Description,
		Info:        event.Info,
		Source:      event.Source,
		Status:      models.ActiveEventStatus,
	}

	if err := r.db.Create(historicalEvent).Error; err != nil {
		return nil, errors.Wrap(err, "failed to create event")
	}
	return historicalEvent, nil
}

func (r *PostgresRepository) UpdateEvent(id int, event *models.UpdateEventDTO) (*models.HistoricalEvent, error) {
	var existingEvent models.HistoricalEvent
	if err := r.db.First(&existingEvent, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, models.ErrEventNotFound
		}
		return nil, errors.Wrap(err, "failed to fetch event")
	}

	if existingEvent.Status == models.DeletedEventStatus {
		return nil, models.ErrEventAlreadyDeleted
	}

	if err := r.db.Model(&existingEvent).Updates(event).Error; err != nil {
		return nil, errors.Wrap(err, "failed to update event")
	}
	return &existingEvent, nil
}

func (r *PostgresRepository) DeleteEvent(id int) error {
	var event models.HistoricalEvent
	if err := r.db.Where("id = ?", id).First(&event).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return models.ErrEventNotFound
		}
		return errors.Wrap(err, "failed to fetch event")
	}

	if event.Status == models.DeletedEventStatus {
		return models.ErrEventAlreadyDeleted
	}

	if event.PhotoURL != nil {
		objectName := strings.TrimPrefix(*event.PhotoURL, fmt.Sprintf("%s/%s/", r.minioClient.EndpointURL().String(), r.bucketName))

		err := r.minioClient.RemoveObject(
			context.Background(),
			r.bucketName,
			objectName,
			minio.RemoveObjectOptions{},
		)
		if err != nil {
			err = errors.Wrap(err, "failed to remove object")
		}
	}

	if err := r.db.Model(&models.HistoricalEvent{}).Where("id = ?", id).Update("status", models.DeletedEventStatus).Update("photo_url", nil).Error; err != nil {
		return errors.Wrap(err, "failed to update event status")
	}

	return nil
}

func (r *PostgresRepository) AddEventToPublication(userID, eventID int) error {
	var event models.HistoricalEvent
	if err := r.db.Where("id = ? AND status = ?", eventID, models.ActiveEventStatus).First(&event).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return models.ErrEventNotFound
		}
		return errors.Wrap(err, "failed to fetch event")
	}

	var publication models.Publication
	err := r.db.Where("status = ? AND user_id = ?", models.DraftPublicationStatus, userID).
		Last(&publication).Error

	if err != nil {
		publication = models.Publication{
			Status: models.DraftPublicationStatus,
			UserID: userID,
		}
		if err := r.db.Create(&publication).Error; err != nil {
			return errors.Wrap(err, "failed to create draft publication")
		}
	}

	publicationEvent := models.PublicationsEvents{
		PublicationID: publication.ID,
		EventID:       eventID,
	}
	if err := r.db.Create(&publicationEvent).Error; err != nil {
		return errors.Wrap(err, "failed to add event to publication")
	}

	return nil
}

func (r *PostgresRepository) AddImageToEvent(id int, file *multipart.FileHeader) error {
	event, err := r.GetEventByID(id)
	if err != nil {
		return err
	}

	objectName := uuid.New().String() + "_" + file.Filename
	fileData, err := file.Open()
	if err != nil {
		return errors.Wrap(err, "failed to open file")
	}
	defer fileData.Close()

	_, err = r.minioClient.PutObject(
		context.Background(),
		r.bucketName,
		objectName,
		fileData,
		file.Size,
		minio.PutObjectOptions{ContentType: file.Header.Get("Content-Type")},
	)
	if err != nil {
		return errors.Wrap(err, "failed to upload file to MinIO")
	}

	photoURL := r.minioClient.EndpointURL().String() + "/" + r.bucketName + "/" + objectName
	event.PhotoURL = &photoURL
	if err := r.db.Save(event).Error; err != nil {
		return errors.Wrap(err, "failed to save photo URL to event")
	}
	return nil
}
