package models

type User struct {
	ID       uint   `gorm:"primaryKey" json:"user_id"`
	FullName string `json:"full_name"`
	Email    string `json:"email"`
	Role     string `json:"role"`
}
