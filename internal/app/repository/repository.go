package repository

import (
	"backend/internal/models"
	"mime/multipart"
	"time"
)

type PostgresRepo interface {
	GetEvents(title string) ([]models.HistoricalEvent, error)
	GetEventByID(id int) (*models.HistoricalEvent, error)
	CreateEvent(event *models.CreateEventDTO) (*models.HistoricalEvent, error)
	UpdateEvent(id int, event *models.UpdateEventDTO) (*models.HistoricalEvent, error)
	DeleteEvent(id int) error
	AddEventToPublication(userID int, eventID int) error
	AddImageToEvent(id int, file *multipart.FileHeader) error
	GetDraftPublication(userID int) (*models.GetPublicationDTO, error)
	GetPublications(userID int, status models.PublicationStatus, startDate, endDate *time.Time) ([]models.Publication, error)
	GetPublicationByID(userID int, id int) (*models.GetPublicationDTO, error)
	UpdatePublication(userID int, id int, publication *models.UpdatePublicationDTO) (*models.GetPublicationDTO, error)
	UpdatePublicationStatus(userID int, id int, status models.PublicationStatus) (*models.GetPublicationDTO, error)
	DeletePublication(userID int, id int) error
	RemoveEventFromPublication(userID int, eventID int) error
	UpdateEventPriority(userID int, publicationID, eventID int, priority int) error
	CreateUser(user *models.RegisterUserDTO) (*models.User, error)
	AuthenticateUser(email, password string) (*models.User, error)
	UpdateUser(userID int, updateData *models.UpdateUserDTO) (*models.User, error)
}

type RedisRepo interface {
	Set(key string, value interface{}, expiration time.Duration) error
	Get(key string) (string, error)
	Delete(key string) error
	Exists(key string) (bool, error)
}
