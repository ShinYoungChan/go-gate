package membership

import (
	"errors"
	"time"
)

type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) GetUserMembership(userID, locationID uint) (*UserMembership, error) {
	userMembership, err := s.repo.GetUserWithMembership(userID, locationID)
	if err != nil {
		return nil, err
	}
	if userMembership == nil {
		return nil, errors.New("회원권을 가지고 있지 않습니다.")
	}
	return userMembership, nil
}

func (s *Service) ValidateEligibility(userID, locationID uint) (*UserMembership, error) {
	userMembership, err := s.repo.GetUserWithMembership(userID, locationID)
	if err != nil {
		return nil, err
	}
	if userMembership == nil {
		return nil, errors.New("회원권을 가지고 있지 않습니다.")
	}

	now := time.Now()
	if now.Before(userMembership.SttDt) || now.After(userMembership.EndDt) {
		return nil, errors.New("이용 가능 기간이 아닙니다.")
	}
	if !userMembership.IsValid {
		return nil, errors.New("정지된 회원권입니다.")
	}
	return userMembership, nil
}

func (s *Service) UpdateMembership(membership *UserMembership) error {
	return s.repo.UpdateUserMembership(membership)
}

func (s *Service) GetTotalAmount(userId uint) (int64, error) {
	return s.repo.SumPaymentAmountByUserID(userId)
}

func (s *Service) GetAvailableMemberships(locationID uint) ([]MembershipItem, error) {
	return s.repo.GetItemsByLocationID(locationID)
}

func (s *Service) GetMembershipItem(itemID uint) (*MembershipItem, error) {
	return s.repo.GetMembershipItem(itemID)
}

func (s *Service) CreateUserMembership(ms *UserMembership) error {
	return s.repo.CreateUserMembership(nil, ms)
}
