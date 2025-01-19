package models

type User struct {
	ID       int    `gorm:"primaryKey" json:"user_id"`
	Name     string `json:"name"`
	Password string `json:"password"`
	Email    string `json:"email"`
	Role     Role   `json:"role"`
}

type Role string

const (
	RoleUser      Role = "user"
	RoleModerator Role = "moderator"
)

type RegisterUserDTO struct {
	Name     string `json:"name" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

type LoginUserDTO struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type UpdateUserDTO struct {
	Name     string `json:"name"`
	Email    string `json:"email" binding:"omitempty,email"`
	Password string `json:"password" binding:"omitempty,min=6"`
}
