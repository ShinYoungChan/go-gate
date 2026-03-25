package repository

import (
	"go-gate/internal/models"

	"gorm.io/gorm"
)

type PaymentRepository struct {
	db *gorm.DB
}

func NewPaymentRepository(db *gorm.DB) *PaymentRepository {
	return &PaymentRepository{db: db}
}

func (r *PaymentRepository) GetDB() *gorm.DB {
	return r.db
}

func (r *PaymentRepository) CreatePaymentLog(log *models.PaymentLog) error {
	return r.db.Create(log).Error
}
