package api

import "time"

type EventType string

const (
	EventTypeLocation EventType = "Локация"
	EventTypeEvent    EventType = "Событие"
	EventTypeArtifact EventType = "Артефакт"
)

type HistoricalEvent struct {
	ID          int       `json:"id"`          // Уникальный идентификатор события
	Title       string    `json:"title"`       // Название события
	EventType   EventType `json:"event_type"`  // тип события: локация, событие, артефакт
	Description string    `json:"description"` // Краткое описание события
	Info        string    `json:"info"`        // информация о соботии
	PhotoURL    string    `json:"photo_url"`   // URL фотографии, связанной с событием
	Source      string    `json:"source"`      // Источник информации об этом событии
}

type Publication struct {
	ID             int        `json:"id"`
	Title          string     `json:"title"`       // Название публикации (исторического события)
	Description    string     `json:"description"` // Краткое описание события
	CreationDate   time.Time  `json:"creation_date"`
	FormationDate  *time.Time `json:"formation_date"`
	CompletionDate *time.Time `json:"completion_date"`
}

type GetPublicationDTO struct {
	Publication
	Events []HistoricalEvent `json:"events"`
}
