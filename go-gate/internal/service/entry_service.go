package service

import (
	"errors"
	"fmt"
	"go-gate/internal/models"
	"go-gate/internal/repository"
	"time"

	"github.com/golang-jwt/jwt/v5"
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

func (s *EntryService) parseAndValidateToken(tokenString string) (uint, uint, error) {
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

	// 주의: JWT 숫자는 기본적으로 float64로 파싱되므로 uint로 형변환 필요
	userID := uint(claims["user_id"].(float64))
	locationID := uint(claims["location_id"].(float64))

	return userID, locationID, nil
}

func (s *EntryService) VerifyEntry(tokenString string, userLat, userLon float64) (*models.UserMembership, error) {
	// 토큰 파싱
	userID, locationID, err := s.parseAndValidateToken(tokenString)

	if err != nil {
		return nil, err // 토큰 만료 or 서명 틀림
	}
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

func (s *EntryService) GenerateEntryToken(userID, locationID uint) (string, error) {
	// 1. 유효한 회원인지 먼저 체크 (이미 짠 로직 재활용)
	_, err := s.membershipService.ValidateEligibility(userID)

	if err != nil {
		return "", err // 회원권 없으면 QR 생성 X
	}
	// 2. Claims 설정
	claims := jwt.MapClaims{
		"user_id":     userID,
		"location_id": locationID,
		"exp":         time.Now().Add(30 * time.Second).Unix(), // 30초 유효
		"iat":         time.Now().Unix(),                       // 토큰발행시간
	}
	// 3. 토큰 생성 및 서명
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// secret_key 하드코딩.. 추후 환경변수로 변경 예정
	secretKey := []byte("qr_secret_key_1234")
	tokenString, err := token.SignedString(secretKey)

	if err != nil {
		return "", errors.New("토큰 생성 실패")
	}
	return tokenString, nil
}
