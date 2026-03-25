package service

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"go-gate/internal/dto"
	"go-gate/internal/models"
	"go-gate/internal/repository"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"gorm.io/gorm"
)

type PaymentService struct {
	paymentRepo    *repository.PaymentRepository
	membershipRepo repository.UserMembershipRepository
}

func NewPaymentService(paymentRepo *repository.PaymentRepository, membershipRepo repository.UserMembershipRepository) *PaymentService {
	return &PaymentService{paymentRepo: paymentRepo, membershipRepo: membershipRepo}
}

func (s *PaymentService) ApprovePayment(req dto.PaymentRequest, userID uint) (*dto.PaymentResponse, error) {
	fmt.Println("회원권 검증 시작")
	// 1. [검증] 기존에 유효한 회원권이 있으면 구매 실패
	userMembership, err := s.membershipRepo.GetUserWithMembership(userID)
	if err != nil {
		return nil, err
	}
	if userMembership != nil && userMembership.IsValid {
		return nil, errors.New("이미 활성화된 회원권이 존재합니다.")

	}
	fmt.Println("회원권 검증 종료")
	// 2. [외부 통신] 토스 API에 '결제 승인' 요청 (http.Post)
	//    - 헤더에 Authorization(Secret Key) 설정
	fmt.Println("토스API  시작")
	secretKey := os.Getenv("TOSS_SECRET_KEY")
	authHeader := "Basic " + base64.StdEncoding.EncodeToString([]byte(secretKey+":"))
	//    - 바디에 paymentKey, orderId, amount 담기
	jsonBody, err := json.Marshal(req)

	if err != nil {
		return nil, err
	}

	tossURL := "https://api.tosspayments.com/v1/payments/confirm"
	httpReq, _ := http.NewRequest("POST", tossURL, bytes.NewBuffer(jsonBody))

	// 헤더 설정
	httpReq.Header.Add("Authorization", authHeader)
	httpReq.Header.Add("Content-Type", "application/json")

	// 실제 전송
	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(httpReq)
	if err != nil {
		return nil, err // 통신 에러
	}
	defer resp.Body.Close()
	fmt.Println("토스API 종료")

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		log.Printf("토스 에러 상세 사유: %s", string(body))
		return nil, errors.New("토스 승인 거절: 결제 정보를 확인해주세요")
	}

	// 3. [성공] DB 트랜잭션 시작
	var res dto.PaymentResponse
	fmt.Println("1. 트랜잭션 진입 직전")
	err = s.paymentRepo.GetDB().Transaction(func(tx *gorm.DB) error {
		fmt.Println("2. 트랜잭션 진입 성공")
		// 1. PaymentLog 저장 (repo.CreatePaymentLog)
		paymentLog := models.PaymentLog{
			UserID:           userID,
			MembershipTypeID: req.MembershipTypeID,
			PaymentKey:       req.PaymentKey,
			OrderId:          req.OrderId,
			Amount:           req.Amount,
			Status:           "DONE",
			Method:           "CARD",
			CreatedAt:        time.Now(),
		}

		if err := s.paymentRepo.CreatePaymentLog(tx, &paymentLog); err != nil {
			return err
		}
		fmt.Println("3. 결제로그 저장 성공")
		// 이용권 없는걸 체크했으니 바로 회원권 생성

		// 이용권 종류가 어떤건지 확인
		membershipType, err := s.membershipRepo.GetMembershipItem(req.MembershipTypeID)

		// 에러가 발생했거나 이용권 종류가 DB에 없는 경우
		if err != nil || membershipType == nil {
			return errors.New("잘못된 이용권 종류입니다.")
		}

		startAt := time.Now()
		endAt := startAt.AddDate(0, 0, membershipType.DurationDays)
		/*  dto 구조체에 시작날짜가 생기면 추가 예정...
		if !req.StartDate.IsZero(){
			startAt = req.StartDate
		}
		*/

		newMembership := models.UserMembership{
			UserID:      userID,
			ItemID:      req.MembershipTypeID,
			SttDt:       startAt,
			EndDt:       endAt,
			IsCountType: membershipType.Type == "count",
			IsValid:     true,
			Amount:      membershipType.Amount,
		}

		if err := s.membershipRepo.CreateUserMembership(tx, &newMembership); err != nil {
			return err
		}
		fmt.Println("4. 이용권 저장 성공")

		res = dto.PaymentResponse{
			Message:   "결제가 완료되었습니다!",
			OrderId:   req.OrderId,
			StartDate: startAt,
			EndDate:   endAt,
		}

		return nil
	})

	log.Printf("트랜잭션 에러 상태: %v", err)
	log.Printf("응답 객체 상태: %+v", res)

	// 트랜잭션 오류 발생시
	if err != nil {
		log.Printf("결제 승인 트랜잭션 실패: %v", err)
		return nil, err
	}

	return &res, nil
}
