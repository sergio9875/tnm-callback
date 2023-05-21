package repository

import "tnm-malawi/connectors/callback/models"

// Repository represent the repositories
type Repository interface {
	Close() error
	GetMbt(mbtId string) (*models.MbtEntity, error)
}
