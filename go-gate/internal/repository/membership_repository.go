package repository

import (
	"errors"
	"go-gate/internal/models"

	"gorm.io/gorm"
)

type MembershipRepository interface {
	// 유저 회원권 관련
	GetUserWithMembership(userID, locationID uint) (*models.UserMembership, error)
	UpdateUserMembership(membership *models.UserMembership) error
	CreateUserMembership(tx *gorm.DB, membership *models.UserMembership) error
	SumPaymentAmountByUserID(userId uint) (int64, error)

	// 멤버십 관련
	GetMembershipItem(itemID uint) (*models.MembershipItem, error)
	GetItemsByLocationID(locationID uint) ([]models.MembershipItem, error)
}

type membershipRepository struct {
	db *gorm.DB
}

func NewUserMembershipRepository(db *gorm.DB) MembershipRepository {
	return &membershipRepository{db: db}
}

func (r *membershipRepository) GetUserWithMembership(userID, locationID uint) (*models.UserMembership, error) {
	var userMembership models.UserMembership
	err := r.db.Where("user_id = ? AND location_id = ?", userID, locationID).First(&userMembership).Error
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

func (r *membershipRepository) UpdateUserMembership(membership *models.UserMembership) error {
	return r.db.Save(membership).Error
}

func (r *membershipRepository) GetMembershipItem(itemID uint) (*models.MembershipItem, error) {
	var membershipItem models.MembershipItem
	err := r.db.Where("id = ?", itemID).First(&membershipItem).Error

	if err != nil {
		return nil, err
	}

	return &membershipItem, nil
}

func (r *membershipRepository) CreateUserMembership(tx *gorm.DB, membership *models.UserMembership) error {
	db := r.db
	if tx != nil {
		db = tx
	}
	return db.Create(membership).Error
}

func (r *membershipRepository) SumPaymentAmountByUserID(userId uint) (int64, error) {
	var total int64
	err := r.db.Model(&models.UserMembership{}).Where("user_id = ?", userId).Select("COALESCE(SUM(amount), 0)").Scan(&total).Error
	return total, err
}

func (r *membershipRepository) GetItemsByLocationID(locationID uint) ([]models.MembershipItem, error) {
	var membershipItems []models.MembershipItem

	err := r.db.Model(&models.MembershipItem{}).Where("location_id = ?", locationID).Find(&membershipItems).Error

	return membershipItems, err
}
