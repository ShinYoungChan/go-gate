package service

import "go-gate/internal/repository"

type AccessLogService struct {
	repo *repository.AccessLogRepository
}

func NewAccessLogService(repo *repository.AccessLogRepository) *AccessLogService {
	return &AccessLogService{repo: repo}
}

func (s *AccessLogService) GetEntryCount(userId uint) (int64, error) {
	return s.repo.CountByUserID(userId)
}
