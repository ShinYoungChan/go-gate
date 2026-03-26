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

func (r *PaymentRepository) CreatePaymentLog(tx *gorm.DB, log *models.PaymentLog) error {
	db := r.db
	if tx != nil {
		db = tx
	}
	return db.Create(log).Error
}
