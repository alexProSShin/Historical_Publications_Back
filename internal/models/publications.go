package models

import "time"

type PublicationStatus string

const (
	DraftPublicationStatus     = "черновик"
	WorkPublicationStatus      = "в работе"
	CompletedPublicationStatus = "завершен"
	RejectedPublicationStatus  = "отклонен"
	DeletedPublicationStatus   = "удален"
)

type UpdatePublicationDTO struct {
	Title       string `json:"title"`
	Description string `json:"description"`
}
type Publication struct {
	ID             int               `gorm:"primarykey" json:"id"`
	Title          string            `gorm:"default:Черновая публикация" json:"title"`                    // Название публикации (исторического события)
	Status         PublicationStatus `gorm:"default:черновик" json:"status"`                              // Установлен статус по умолчанию
	Description    string            `gorm:"default:Это описание черновой публикации" json:"description"` // Краткое описание события
	TrustScore     *float64          `gorm:"default:0" json:"trust_score"`
	CreationDate   time.Time         `gorm:"default:CURRENT_TIMESTAMP" json:"creation_date"`
	FormationDate  *time.Time        `json:"formation_date"`
	CompletionDate *time.Time        `json:"completion_date"`
	UserID         int               `json:"user_id"`
	ModeratorID    *int              `json:"moderator_id"`
}

type GetPublicationDTO struct {
	Publication
	Events []HistoricalEvent `json:"events"`
}
