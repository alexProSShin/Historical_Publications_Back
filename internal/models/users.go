package models

type User struct {
	ID       int    `gorm:"primaryKey" json:"user_id" binding:"required"`
	Name     string `json:"name" binding:"required"`
	Password string `json:"password" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Role     Role   `json:"role" binding:"required"`
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

type LoginResponseDTO struct {
	Token string `json:"token" binding:"required"`
	User  User   `json:"user" binding:"required"`
}

type UpdateUserDTO struct {
	Name     string `json:"name" binding:"omitempty"`
	Email    string `json:"email" binding:"omitempty,email"`
	Password string `json:"password" binding:"omitempty,min=6"`
}
