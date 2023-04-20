package repository

import "malawi-callback/models"

// Repository represent the repositories
type Repository interface {
	Close() error
	GetMbt(mbtId string) (*models.MbtEntity, error)
}
