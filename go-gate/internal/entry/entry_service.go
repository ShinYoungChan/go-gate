package entry

import (
	"errors"
	"fmt"
	"go-gate/internal/location"
	"go-gate/internal/membership"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type Service struct {
	membershipService *membership.Service
	logRepo           *AccessLogRepository
	locationService   *location.Service
}

func NewService(ms *membership.Service, logRepo *AccessLogRepository, locService *location.Service) *Service {
	return &Service{
		membershipService: ms,
		logRepo:           logRepo,
		locationService:   locService,
	}
}

func (s *Service) parseAndValidateToken(tokenString string) (uint, uint, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte("qr_secret_key_1234"), nil
	})
	if err != nil || !token.Valid {
		return 0, 0, errors.New("유효하지 않거나 만료된 QR 코드입니다.")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return 0, 0, errors.New("토큰 데이터 형식이 잘못되었습니다.")
	}

	userID := uint(claims["user_id"].(float64))
	locationID := uint(claims["location_id"].(float64))
	return userID, locationID, nil
}

func (s *Service) VerifyEntry(tokenString string, userLat, userLon float64) (*membership.UserMembership, error) {
	userID, locationID, err := s.parseAndValidateToken(tokenString)
	if err != nil {
		return nil, err
	}

	loc, err := s.locationService.GetLocation(locationID)
	if err != nil {
		return nil, err
	}

	distance := s.locationService.CalculateDistance(userLat, userLon, loc.Lat, loc.Lon)
	fmt.Printf("사용자: %f, %f / DB장소: %f, %f\n", userLat, userLon, loc.Lat, loc.Lon)
	if distance > 50 {
		fmt.Println("Distance = ", distance)
		return nil, errors.New("장소와 너무 멉니다. 입구 근처에서 다시 시도해주세요.")
	}

	userMembership, err := s.membershipService.ValidateEligibility(userID, locationID)
	if err != nil {
		return nil, err
	}

	lastLog, err := s.logRepo.GetLastAccessLog(userID)
	if lastLog != nil {
		fmt.Println("최근 입장 시간:", lastLog.AccessedAt)
		if time.Since(lastLog.AccessedAt) < 1*time.Minute {
			return nil, errors.New("방금 입장하셨습니다. 잠시 후 다시 시도해주세요.")
		}
	}

	if userMembership.IsCountType {
		if userMembership.Count > 0 {
			userMembership.Count--
		} else {
			return nil, errors.New("횟수권을 모두 사용했습니다.")
		}
	}

	if err = s.membershipService.UpdateMembership(userMembership); err != nil {
		return nil, errors.New("업데이트 실패")
	}

	logData := AccessLog{
		UserID:     userID,
		AccessedAt: time.Now(),
		Method:     "QR코드",
		Result:     "성공",
	}
	if err = s.logRepo.CreateEntryLog(&logData); err != nil {
		return nil, errors.New("로그 저장 실패")
	}

	return userMembership, nil
}

func (s *Service) GenerateEntryToken(userID, locationID uint) (string, error) {
	_, err := s.membershipService.ValidateEligibility(userID, locationID)
	if err != nil {
		return "", err
	}

	claims := jwt.MapClaims{
		"user_id":     userID,
		"location_id": locationID,
		"exp":         time.Now().Add(30 * time.Second).Unix(),
		"iat":         time.Now().Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// secret_key 하드코딩 — 추후 환경변수로 변경 예정
	tokenString, err := token.SignedString([]byte("qr_secret_key_1234"))
	if err != nil {
		return "", errors.New("토큰 생성 실패")
	}
	return tokenString, nil
}
