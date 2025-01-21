package models

type PublicationsEvents struct {
	PublicationID int `gorm:"type:serial;primaryKey;index" json:"publication_id" binding:"required"`
	EventID       int `gorm:"type:serial;primaryKey;index" json:"event_id" binding:"required"`
	Priority      int `gorm:"default:0" json:"priority" binding:"required"`
}

type UpdateEventPriority struct {
	Priority int `json:"priority"`
}
