package repository

import (
	"errors"
	"go-gate/internal/models"

	"gorm.io/gorm"
)

type UserMembershipRepository struct {
	db *gorm.DB
}

func NewUserMembershipRepository(db *gorm.DB) *UserMembershipRepository {
	return &UserMembershipRepository{db: db}
}

func (r *UserMembershipRepository) GetUserWithMembership(userID uint) (*models.UserMembership, error) {
	var userMembership models.UserMembership
	err := r.db.Where("user_id = ?", userID).First(&userMembership).Error
	if err != nil {
		// 에러가 데이터를 찾지 못한 건지 확인
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		// 그 외는 실제 에러(DB 연결 끊김)등
		return nil, err
	}
	return &userMembership, nil
}

func (r *UserMembershipRepository) UpdateMembership(membership *models.UserMembership) error {
	return r.db.Save(membership).Error
}
