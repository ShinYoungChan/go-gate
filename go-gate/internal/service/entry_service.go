package service

import (
	"errors"
	"fmt"
	"go-gate/internal/models"
	"go-gate/internal/repository"
	"time"
)

type EntryService struct {
	membershipService *UserMembershipService
	logRepo           *repository.AccessLogRepository
	locationService   *LocationService
}

func NewEntryService(membershipService *UserMembershipService, logRepo *repository.AccessLogRepository, locService *LocationService) *EntryService {
	return &EntryService{
		membershipService: membershipService,
		logRepo:           logRepo,
		locationService:   locService,
	}
}

func (s *EntryService) VerifyEntry(userID uint, userLat, userLon float64, locationID uint) (*models.UserMembership, error) {
	// 1. DB에서 장소(Location) 정보 가져오기
	location, err := s.locationService.GetLocation(locationID)
	if err != nil {
		return nil, err
	}

	// 2. 거리 계산 호출!
	distance := s.locationService.CalculateDistance(userLat, userLon, location.Lat, location.Lon)

	fmt.Printf("사용자: %f, %f / DB장소: %f, %f\n", userLat, userLon, location.Lat, location.Lon)
	// 3. 거리 체크 (예: 50m 이내)
	if distance > 50 {
		fmt.Println("Distance = ", distance)
		return nil, errors.New("장소와 너무 멉니다. 입구 근처에서 다시 시도해주세요.")
	}

	// 4. 여기서부터는 멤버십 체크 로직 (SttDt, EndAt, Count 등)
	userMembership, err := s.membershipService.ValidateEligibility(userID)

	if err != nil {
		return nil, err
	}

	// 5. 1분안에 입장했던 기록이있는지
	lastLog, err := s.logRepo.GetLastAccessLog(userID)
	if lastLog != nil {
		fmt.Println("최근 입장 시간:", lastLog.AccessedAt)
		if time.Since(lastLog.AccessedAt) < 1*time.Minute {
			return nil, errors.New("방금 입장하셨습니다. 잠시 후 다시 시도해주세요.")
		}
	}
	// 6. 차감 및 저장, 입장 로그 저장 access_log 호출
	err = s.membershipService.UpdateMembership(userMembership)

	if err != nil {
		// 오류문구 추후 수정
		return nil, errors.New("업데이트 실패")
	}

	logData := models.AccessLog{
		UserID:     userID,
		AccessedAt: time.Now(),
		Method:     "QR코드",
		Result:     "성공",
	}

	err = s.logRepo.CreateEntryLog(&logData)

	if err != nil {
		// 추후 문구 수정
		return nil, errors.New("로그 저장 실패")
	}

	return userMembership, nil
}

func (s *EntryService) GenerateEntryToken(userID uint) (string, error) {
	// 1. 유효한 회원인지 먼저 체크 (이미 짠 로직 재활용)
	_, err := s.membershipService.ValidateEligibility(userID)

	if err != nil {
		return "", err // 회원권 없으면 QR 생성 X
	}
	// 2. 현재 시간과 정보를 섞어서 암호화된 문자열 생성
	// 3. Redis나 메모리에 "이 토큰은 30초간 유효해"라고 저장 (선택사항)
	return "token_abc_123", nil
}
