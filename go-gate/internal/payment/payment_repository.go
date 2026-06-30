package payment

import (
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

func (r *PaymentRepository) CreatePaymentLog(tx *gorm.DB, log *PaymentLog) error {
	db := r.db
	if tx != nil {
		db = tx
	}
	return db.Create(log).Error
}

func (r *PaymentRepository) PaymentLogFindByUserID(userID uint) ([]PaymentLog, error) {
	var paymentLogs []PaymentLog
	err := r.db.Where("user_id = ?", userID).Order("created_at DESC").Find(&paymentLogs).Error
	return paymentLogs, err
}
