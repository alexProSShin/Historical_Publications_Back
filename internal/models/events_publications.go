package models

type PublicationsEvents struct {
	PublicationID int `gorm:"type:serial;primaryKey;index" json:"publication_id"`
	EventID       int `gorm:"type:serial;primaryKey;index" json:"event_id"`
	Priority      int `gorm:"default:0" json:"priority"`
}

type UpdateEventPriority struct {
	Priority int `json:"priority"`
}
