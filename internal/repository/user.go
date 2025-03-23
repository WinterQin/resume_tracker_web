package repository

import (
	"internship-manager/internal/model"
	"internship-manager/pkg/database"
)

// internal/repository/user.go
func CreateUser(user model.User) error {
	return database.DB.Create(&user).Error
}
