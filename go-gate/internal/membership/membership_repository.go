package membership

import (
	"errors"

	"gorm.io/gorm"
)

type Repository interface {
	GetUserWithMembership(userID, locationID uint) (*UserMembership, error)
	UpdateUserMembership(membership *UserMembership) error
	CreateUserMembership(tx *gorm.DB, membership *UserMembership) error
	SumPaymentAmountByUserID(userId uint) (int64, error)
	GetMembershipItem(itemID uint) (*MembershipItem, error)
	GetItemsByLocationID(locationID uint) ([]MembershipItem, error)
}

type membershipRepository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &membershipRepository{db: db}
}

func (r *membershipRepository) GetUserWithMembership(userID, locationID uint) (*UserMembership, error) {
	var userMembership UserMembership
	err := r.db.Where("user_id = ? AND location_id = ?", userID, locationID).First(&userMembership).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &userMembership, nil
}

func (r *membershipRepository) UpdateUserMembership(membership *UserMembership) error {
	return r.db.Save(membership).Error
}

func (r *membershipRepository) GetMembershipItem(itemID uint) (*MembershipItem, error) {
	var membershipItem MembershipItem
	err := r.db.Where("id = ?", itemID).First(&membershipItem).Error
	if err != nil {
		return nil, err
	}
	return &membershipItem, nil
}

func (r *membershipRepository) CreateUserMembership(tx *gorm.DB, membership *UserMembership) error {
	db := r.db
	if tx != nil {
		db = tx
	}
	return db.Create(membership).Error
}

func (r *membershipRepository) SumPaymentAmountByUserID(userId uint) (int64, error) {
	var total int64
	err := r.db.Model(&UserMembership{}).Where("user_id = ?", userId).Select("COALESCE(SUM(amount), 0)").Scan(&total).Error
	return total, err
}

func (r *membershipRepository) GetItemsByLocationID(locationID uint) ([]MembershipItem, error) {
	var membershipItems []MembershipItem
	err := r.db.Model(&MembershipItem{}).Where("location_id = ?", locationID).Find(&membershipItems).Error
	return membershipItems, err
}
