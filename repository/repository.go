package repository

import "malawi-callback/models"

// Repository represent the repositories
type Repository interface {
  Close() error
  FindUserByID(id int) (*models.UserEntity, error)
  CreateUser(user *models.UserEntity) error
  UpdateUser(user *models.UserEntity) error
}
