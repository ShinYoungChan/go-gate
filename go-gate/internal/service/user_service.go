package service

import (
	"errors"
	"fmt"
	"go-gate/internal/models"
	"go-gate/internal/repository"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	repo            *repository.UserRepository
	membershipRepo  *repository.UserMembershipRepository
	logRepo         *repository.AccessLogRepository
	locationService *LocationService
}

func NewUserService(repo *repository.UserRepository, membershipRepo *repository.UserMembershipRepository, logRepo *repository.AccessLogRepository, locService *LocationService) *UserService {
	return &UserService{
		repo:            repo,
		membershipRepo:  membershipRepo,
		logRepo:         logRepo,
		locationService: locService,
	}
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

func (s *UserService) VerifyEntry(userID uint, userLat, userLon float64, locationID uint) error {
	// 1. DB에서 장소(Location) 정보 가져오기
	// (이 부분은 LocationRepo를 통해 가져와야겠죠?)
	location, err := s.locationService.GetLocation(locationID)
	if err != nil {
		return err
	}

	// 2. 거리 계산 호출!
	distance := s.locationService.CalculateDistance(userLat, userLon, location.Lat, location.Lon)

	fmt.Printf("사용자: %f, %f / DB장소: %f, %f\n", userLat, userLon, location.Lat, location.Lon)
	// 3. 거리 체크 (예: 50m 이내)
	if distance > 50 {
		fmt.Println("Distance = ", distance)
		return errors.New("장소와 너무 멉니다. 입구 근처에서 다시 시도해주세요.")
	}

	// 4. 여기서부터는 멤버십 체크 로직 (SttDt, EndAt, Count 등)
	userMembership, err := s.membershipRepo.GetUserWithMembership(userID)
	if userMembership == nil {
		// 일단 에러 리턴, 추후 회원권 구매로직으로 추가 예정..?
		return errors.New("회원권을 가지고 있지 않습니다.")
	}
	// 4-1. 기간 체크
	now := time.Now()
	// 시작기간 전 이거나 종료기한 이후면 에러 처리
	if now.Before(userMembership.SttDt) || now.After(userMembership.EndDt) {
		return errors.New("이용 가능 기간이 아닙니다.")
	}

	// 4-1-1. 회원권 사용 여부 체크
	if !userMembership.IsValid {
		return errors.New("정지된 회원권입니다.")
	}

	// 4-1-2. 1분안에 입장했던 기록이있는지
	lastLog, err := s.logRepo.GetLastAccessLog(userID)
	if lastLog != nil {
		fmt.Println("최근 입장 시간:", lastLog.AccessedAt)
		if time.Since(lastLog.AccessedAt) < 1*time.Minute {
			return errors.New("방금 입장하셨습니다. 잠시 후 다시 시도해주세요.")
		}
	}

	// 4-2. 회원권 종류 체크(정기권, 횟수권)
	if userMembership.IsCountType {
		// 4-3. 횟수권인 경우 잔여 횟수 확인
		if userMembership.Count > 0 {
			userMembership.Count--
		} else {
			// 우선 에러 리턴, 이후 회원군 구매 로직으로 이동..
			return errors.New("횟수권을 모두 사용했습니다.")
		}
	}
	// 4-4. 차감 및 저장, 입장 로그 저장 access_log 호출
	err = s.membershipRepo.UpdateMembership(userMembership)

	if err != nil {
		// 오류문구 추후 수정
		return errors.New("업데이트 실패")
	}

	logData := models.AccessLog{
		UserID:     userID,
		AccessedAt: now,
		Method:     "QR코드",
		Result:     "성공",
	}

	err = s.logRepo.CreateEntryLog(&logData)

	if err != nil {
		// 추후 문구 수정
		return errors.New("로그 저장 실패")
	}

	return nil
}
