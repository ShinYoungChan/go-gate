package service

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"errors"
	"go-gate/internal/dto"
	"go-gate/internal/models"
	"go-gate/internal/repository"
	"net/http"
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

func (s *PaymentService) ApprovePayment(req dto.PaymentRequest, userID uint) error {
	// 1. [외부 통신] 토스 API에 '결제 승인' 요청 (http.Post)
	//    - 헤더에 Authorization(Secret Key) 설정
	secretKey := "payment_secret_key_1234"
	authHeader := "Basic" + base64.StdEncoding.EncodeToString([]byte(secretKey+":"))
	//    - 바디에 paymentKey, orderId, amount 담기
	jsonBody, err := json.Marshal(req)

	if err != nil {
		return err
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
		return err // 통신 에러
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return errors.New("토스 승인 거절: 결제 정보를 확인해주세요")
	}
	// [성공] DB 트랜잭션 시작
	return s.paymentRepo.GetDB().Transaction(func(tx *gorm.DB) error {
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

		if err := s.paymentRepo.CreatePaymentLog(&paymentLog); err != nil {
			return err
		}
		// 2. UserMembership 업데이트 (membershipRepo.Update)
		// 가지고있는 회원권이 있는지 체크
		userMembership, err := s.membershipRepo.GetUserWithMembership(userID)

		if err != nil {
			return err
		}

		return nil
	})
}
