package repository

import (
	"errors"
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

func (r *AccessLogRepository) GetLastAccessLog(userID uint) (*models.AccessLog, error) {
	var log models.AccessLog

	query := r.db.Where("user_id = ?", userID)
	query = query.Order("accessed_at DESC")
	query = query.Limit(1)

	err := query.First(&log).Error
	if err != nil {
		// 에러가 데이터를 찾지 못한 건지 확인
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		// 그 외는 실제 에러(DB 연결 끊김)등
		return nil, err
	}

	return &log, nil
}
