package models

type EventType string

const (
	EventTypeLocation EventType = "Локация"
	EventTypeEvent    EventType = "Событие"
	EventTypeArtifact EventType = "Артефакт"
)

type EventStatus string

const (
	ActiveEventStatus  EventStatus = "активно"
	DeletedEventStatus EventStatus = "удалено"
)

type CreateEventDTO struct {
	Title       string    `json:"title" binding:"required"`       // Название события
	EventType   EventType `json:"event_type" binding:"required"`  // тип события: локация, событие, артефакт
	Description string    `json:"description" binding:"required"` // Краткое описание события
	Info        string    `json:"info" binding:"required"`        // информация о событии
	Source      *string   `json:"source" binding:"omitempty"`     // Источник информации об этом событии
}

type UpdateEventDTO struct {
	Title       string    `json:"title" binding:"omitempty"`       // Название события
	EventType   EventType `json:"event_type" binding:"omitempty"`  // тип события: локация, событие, артефакт
	Description string    `json:"description" binding:"omitempty"` // Краткое описание события
	Info        string    `json:"info" binding:"omitempty"`        // информация о событии
	Source      *string   `json:"source" binding:"omitempty"`      // Источник информации об этом событии
}

type GetEvents struct {
	HistoricalEvent
	Priority *string `json:"priority,omitempty" binding:"omitempty"` // Приоритет
}

type HistoricalEvent struct {
	ID          int         `gorm:"primarykey" json:"id" binding:"required"`          // Уникальный идентификатор события
	Status      EventStatus `gorm:"default:активно" json:"status" binding:"required"` // Установлен статус по умолчанию
	Title       string      `json:"title" binding:"required"`                         // Название события
	EventType   EventType   `json:"event_type" binding:"required"`                    // тип события: локация, событие, артефакт
	Description string      `json:"description" binding:"required"`                   // Краткое описание события
	Info        string      `json:"info" binding:"required"`                          // информация о событии
	PhotoURL    *string     `json:"photo_url" binding:"omitempty"`                    // URL фотографии, связанной с событием
	Source      *string     `json:"source" binding:"omitempty"`                       // Источник информации об этом событии
}

type GetEventsDTO struct {
	HistoricalEvents []HistoricalEvent `json:"historical_events" binding:"omitempty"`
	PublicationID    int               `json:"publications_id" binding:"required"`
	EventsCount      int               `json:"events_count" binding:"required"`
}
