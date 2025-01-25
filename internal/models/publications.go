package models

import "time"

type PublicationStatus string

const (
	DraftPublicationStatus     PublicationStatus = "черновик"
	WorkPublicationStatus      PublicationStatus = "в работе"
	CompletedPublicationStatus PublicationStatus = "завершен"
	RejectedPublicationStatus  PublicationStatus = "отклонен"
	DeletedPublicationStatus   PublicationStatus = "удален"
)

type UpdatePublicationDTO struct {
	Title       string `json:"title" binding:"required"`
	Description string `json:"description" binding:"required"`
}

type GetPublications struct {
	Publication
	UserName *string `json:"user_name,omitempty" binding:"omitempty"`
}

type Publication struct {
	ID             int               `gorm:"primarykey" json:"id" binding:"required"`
	Title          string            `gorm:"default:Черновая публикация" json:"title" binding:"required"`                    // Название публикации (исторического события)
	Status         PublicationStatus `gorm:"default:черновик" json:"status" binding:"required"`                              // Установлен статус по умолчанию
	Description    string            `gorm:"default:Это описание черновой публикации" json:"description" binding:"required"` // Краткое описание события
	TrustScore     *float64          `gorm:"default:0" json:"trust_score" binding:"omitempty"`
	CreationDate   time.Time         `gorm:"default:CURRENT_TIMESTAMP" json:"creation_date" binding:"required"`
	FormationDate  *time.Time        `json:"formation_date" binding:"omitempty"`
	CompletionDate *time.Time        `json:"completion_date" binding:"omitempty"`
	UserID         int               `json:"user_id" binding:"required"`
	ModeratorID    *int              `json:"moderator_id" binding:"omitempty"`
}

type GetPublicationDTO struct {
	Publication
	Events []HistoricalEvent `json:"events" binding:"omitempty"`
}
