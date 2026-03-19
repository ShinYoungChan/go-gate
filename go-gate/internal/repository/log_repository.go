package repository

import (
	"go-gate/internal/models"

	"gorm.io/gorm"
)

type AccessLogRepository struct {
	db *gorm.DB
}

func NewAccessLogRepository(db *gorm.DB) *AccessLogRepository {
	return &AccessLogRepository{db: db}
}

func (r *AccessLogRepository) CreateEntryLog(log *models.AccessLog) error {
	return r.db.Create(log).Error
}
