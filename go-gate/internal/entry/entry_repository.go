package entry

import (
	"errors"

	"gorm.io/gorm"
)

type AccessLogRepository struct {
	db *gorm.DB
}

func NewAccessLogRepository(db *gorm.DB) *AccessLogRepository {
	return &AccessLogRepository{db: db}
}

func (r *AccessLogRepository) CreateEntryLog(log *AccessLog) error {
	return r.db.Create(log).Error
}

func (r *AccessLogRepository) GetLastAccessLog(userID uint) (*AccessLog, error) {
	var log AccessLog
	err := r.db.Where("user_id = ?", userID).Order("accessed_at DESC").Limit(1).First(&log).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &log, nil
}

func (r *AccessLogRepository) CountByUserID(userId uint) (int64, error) {
	var count int64
	err := r.db.Model(&AccessLog{}).Where("user_id = ?", userId).Count(&count).Error
	return count, err
}
