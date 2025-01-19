package postgresrepo

import (
	"backend/internal/models"
	"golang.org/x/crypto/bcrypt"
)

func (r *PostgresRepository) CreateUser(user *models.RegisterUserDTO) (*models.User, error) {
	existingUser := &models.User{}
	if err := r.db.Where("email = ?", user.Email).First(existingUser).Error; err == nil {
		return nil, models.ErrUserAlreadyExists
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	newUser := &models.User{
		Name:     user.Name,
		Email:    user.Email,
		Password: string(hashedPassword),
		Role:     models.RoleUser,
	}
	if err := r.db.Create(newUser).Error; err != nil {
		return nil, err
	}

	return newUser, nil
}

func (r *PostgresRepository) AuthenticateUser(email, password string) (*models.User, error) {
	user := &models.User{}
	if err := r.db.Where("email = ?", email).First(user).Error; err != nil {
		return nil, models.ErrInvalidCredentials
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return nil, models.ErrInvalidCredentials
	}

	return user, nil
}

func (r *PostgresRepository) UpdateUser(userID int, updateData *models.UpdateUserDTO) (*models.User, error) {
	user := &models.User{}
	if err := r.db.First(user, userID).Error; err != nil {
		return nil, models.ErrUserNotFound
	}

	currentRole := user.Role
	if updateData.Name != "" {
		user.Name = updateData.Name
	}
	if updateData.Email != "" {
		user.Email = updateData.Email
	}
	if updateData.Password != "" {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(updateData.Password), bcrypt.DefaultCost)
		if err != nil {
			return nil, err
		}
		user.Password = string(hashedPassword)
	}
	user.Role = currentRole

	if err := r.db.Save(user).Error; err != nil {
		return nil, err
	}

	return user, nil
}

func (r *PostgresRepository) GetUserByID(userID int) (*models.User, error) {
	var user models.User
	if err := r.db.First(&user, userID).Error; err != nil {
		return nil, err
	}
	return &user, nil
}
