package service

import (
	"errors"
	"go-gate/internal/models"
	"go-gate/internal/repository"
	"time"
)

type UserMembershipService struct {
	repo repository.UserMembershipRepository
}

func NewUserMembershipService(repo repository.UserMembershipRepository) *UserMembershipService {
	return &UserMembershipService{repo: repo}
}

func (s *UserMembershipService) ValidateEligibility(userID uint) (*models.UserMembership, error) {
	userMembership, err := s.repo.GetUserWithMembership(userID)
	if err != nil {
		return nil, err
	}

	if userMembership == nil {
		// 일단 에러 리턴, 추후 회원권 구매로직으로 추가 예정..?
		return nil, errors.New("회원권을 가지고 있지 않습니다.")
	}
	// 1. 기간 체크
	now := time.Now()
	// 시작기간 전 이거나 종료기한 이후면 에러 처리
	if now.Before(userMembership.SttDt) || now.After(userMembership.EndDt) {
		return nil, errors.New("이용 가능 기간이 아닙니다.")
	}

	// 2. 회원권 사용 여부 체크
	if !userMembership.IsValid {
		return nil, errors.New("정지된 회원권입니다.")
	}

	// 3. 회원권 종류 체크(정기권, 횟수권)
	if userMembership.IsCountType {
		// 3-1. 횟수권인 경우 잔여 횟수 확인
		if userMembership.Count > 0 {
			userMembership.Count--
		} else {
			// 우선 에러 리턴, 이후 회원군 구매 로직으로 이동..
			return nil, errors.New("횟수권을 모두 사용했습니다.")
		}
	}
	return userMembership, nil
}

func (s *UserMembershipService) UpdateMembership(membership *models.UserMembership) error {
	return s.repo.UpdateUserMembership(membership)
}
