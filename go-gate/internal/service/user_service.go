package service

import (
	"errors"
	"go-gate/internal/models"
	"go-gate/internal/repository"

	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	repo *repository.UserRepository
}

func NewUserService(repo *repository.UserRepository) *UserService {
	return &UserService{repo: repo}
}

func (s *UserService) SignUpUser(name, email, password string) error {
	// 이메일 중복체크
	existingUser, err := s.repo.GetUserByEmail(email)
	if err != nil {
		return err // DB 연결 오류 등의 에러 처리
	}

	if existingUser != nil {
		return errors.New("already exists")
	}

	// 비밀번호 암호화
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	if err != nil {
		return err
	}

	// 모델생성
	user := models.User{
		Name:     name,
		Email:    email,
		Password: string(hashedPassword),
	}

	// repo를 통해 DB 저장
	return s.repo.CreateUser(&user)
}

func (s *UserService) AuthenticateUser(email, password string) error {
	existingUser, err := s.repo.GetUserByEmail(email)
	if err != nil {
		return err // DB 연결 오류 등의 에러 처리
	}

	if existingUser == nil {
		return errors.New("아이디 혹은 비밀번호가 일치하지 않습니다.")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(existingUser.Password), []byte(password)); err != nil {
		return errors.New("아이디 혹은 비밀번호가 일치하지 않습니다.")
	}
	return nil
}

func (s *UserService) GetUser(userID uint) (*models.User, error) {
	user, err := s.repo.GetUserById(userID)
	if err != nil {
		return nil, err
	}

	if user == nil {
		return nil, errors.New("회원권을 가지고 있지 않습니다.")
	}

	return user, nil
}
