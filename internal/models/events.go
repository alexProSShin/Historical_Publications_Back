package models

type EventType string

const (
	EventTypeLocation EventType = "Локация"
	EventTypeEvent    EventType = "Событие"
	EventTypeArtifact EventType = "Артефакт"
)

type EventStatus string

const (
	ActiveEventStatus  = "активно"
	DeletedEventStatus = "удалено"
)

type HistoricalEvent struct {
	ID          int         `gorm:"primarykey" json:"id"`          // Уникальный идентификатор события
	Status      EventStatus `gorm:"default:активно" json:"status"` // Установлен статус по умолчанию
	Title       string      `json:"title"`                         // Название события
	EventType   EventType   `json:"event_type"`                    // тип события: локация, событие, артефакт
	Description string      `json:"description"`                   // Краткое описание события
	Info        string      `json:"info"`                          // информация о событии
	PhotoURL    *string     `json:"photo_url"`                     // URL фотографии, связанной с событием
	Source      *string     `json:"source"`                        // Источник информации об этом событии
}

type GetEventsDTO struct {
	HistoricalEvents []HistoricalEvent `json:"historical_events"`
	PublicationID    int               `json:"publications_id"`
}
