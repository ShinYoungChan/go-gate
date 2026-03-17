package service

import (
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
