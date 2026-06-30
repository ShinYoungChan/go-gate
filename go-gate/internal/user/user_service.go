package user

import (
	"errors"
	"go-gate/internal/entry"
	"go-gate/internal/membership"

	"golang.org/x/crypto/bcrypt"
)

type SummaryResponse struct {
	EntryCount  int64 `json:"entry_count"`
	TotalAmount int64 `json:"total_amount"`
}

type Service struct {
	repo              *UserRepository
	accessLogRepo     *entry.AccessLogRepository
	membershipService *membership.Service
}

func NewService(repo *UserRepository, alRepo *entry.AccessLogRepository, ms *membership.Service) *Service {
	return &Service{
		repo:              repo,
		accessLogRepo:     alRepo,
		membershipService: ms,
	}
}

func (s *Service) SignUpUser(name, email, password string) error {
	existingUser, err := s.repo.GetUserByEmail(email)
	if err != nil {
		return err
	}
	if existingUser != nil {
		return errors.New("already exists")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	return s.repo.CreateUser(&User{
		Name:     name,
		Email:    email,
		Password: string(hashedPassword),
	})
}

func (s *Service) AuthenticateUser(email, password string) error {
	existingUser, err := s.repo.GetUserByEmail(email)
	if err != nil {
		return err
	}
	if existingUser == nil {
		return errors.New("아이디 혹은 비밀번호가 일치하지 않습니다.")
	}
	if err := bcrypt.CompareHashAndPassword([]byte(existingUser.Password), []byte(password)); err != nil {
		return errors.New("아이디 혹은 비밀번호가 일치하지 않습니다.")
	}
	return nil
}

func (s *Service) GetUser(userID uint) (*User, error) {
	user, err := s.repo.GetUserById(userID)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errors.New("회원권을 가지고 있지 않습니다.")
	}
	return user, nil
}

func (s *Service) GetUserSummary(userId uint) (SummaryResponse, error) {
	count, err := s.accessLogRepo.CountByUserID(userId)
	if err != nil {
		return SummaryResponse{}, err
	}

	total, err := s.membershipService.GetTotalAmount(userId)
	if err != nil {
		return SummaryResponse{}, err
	}

	return SummaryResponse{
		EntryCount:  count,
		TotalAmount: total,
	}, nil
}
