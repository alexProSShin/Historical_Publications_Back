package models

type PublicationsEvents struct {
	PublicationID int `gorm:"type:serial;primaryKey;index" json:"publication_id"`
	EventID       int `gorm:"type:serial;primaryKey;index" json:"event_id"`
}
