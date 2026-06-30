package payment

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"go-gate/internal/membership"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"gorm.io/gorm"
)

type Service struct {
	paymentRepo    *PaymentRepository
	membershipRepo membership.Repository
}

func NewService(paymentRepo *PaymentRepository, membershipRepo membership.Repository) *Service {
	return &Service{paymentRepo: paymentRepo, membershipRepo: membershipRepo}
}

func (s *Service) ApprovePayment(req Request, userID, locationID uint) (*Response, error) {
	fmt.Println("회원권 검증 시작")
	userMembership, err := s.membershipRepo.GetUserWithMembership(userID, locationID)
	if err != nil {
		return nil, err
	}
	if userMembership != nil && userMembership.IsValid {
		return nil, errors.New("이미 활성화된 회원권이 존재합니다.")
	}
	fmt.Println("회원권 검증 종료")

	fmt.Println("토스API 시작")
	secretKey := os.Getenv("TOSS_SECRET_KEY")
	authHeader := "Basic " + base64.StdEncoding.EncodeToString([]byte(secretKey+":"))

	jsonBody, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}

	tossURL := "https://api.tosspayments.com/v1/payments/confirm"
	httpReq, _ := http.NewRequest("POST", tossURL, bytes.NewBuffer(jsonBody))
	httpReq.Header.Add("Authorization", authHeader)
	httpReq.Header.Add("Content-Type", "application/json")

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(httpReq)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	fmt.Println("토스API 종료")

	bodyBytes, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusOK {
		log.Printf("토스 에러 상세 사유: %s", string(bodyBytes))
		return nil, errors.New("토스 승인 거절: 결제 정보를 확인해주세요")
	}

	var tossResp map[string]interface{}
	if err := json.Unmarshal(bodyBytes, &tossResp); err != nil {
		return nil, err
	}

	actualAmount := uint64(tossResp["totalAmount"].(float64))
	if actualAmount != uint64(req.Amount) {
		log.Printf("[보안 경고] 금액 불일치! 요청: %d, 실제: %d", req.Amount, actualAmount)
		return nil, errors.New("결제 금액이 요청 정보와 일치하지 않습니다.")
	}

	var res Response
	fmt.Println("1. 트랜잭션 진입 직전")
	err = s.paymentRepo.GetDB().Transaction(func(tx *gorm.DB) error {
		fmt.Println("2. 트랜잭션 진입 성공")
		paymentLog := PaymentLog{
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

		membershipType, err := s.membershipRepo.GetMembershipItem(req.MembershipTypeID)
		if err != nil || membershipType == nil {
			return errors.New("잘못된 이용권 종류입니다.")
		}

		startAt := time.Now()
		endAt := startAt.AddDate(0, 0, membershipType.DurationDays)

		newMembership := membership.UserMembership{
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

		res = Response{
			Message:   "결제가 완료되었습니다!",
			OrderId:   req.OrderId,
			StartDate: startAt,
			EndDate:   endAt,
		}
		return nil
	})

	log.Printf("트랜잭션 에러 상태: %v", err)
	log.Printf("응답 객체 상태: %+v", res)

	if err != nil {
		log.Printf("결제 승인 트랜잭션 실패 결제 취소 시도: %v", err)
		cancelErr := s.cancelTossPayment(req.PaymentKey, "서버 내부 DB 오류로 인한 자동 취소")
		if cancelErr != nil {
			log.Printf("[ERROR] 결제 취소 실패. 수동 환불 필요: %s", req.PaymentKey)
		}
		return nil, errors.New("서비스 처리 중 오류가 발생하여 결제가 자동 취소되었습니다.")
	}

	return &res, nil
}

func (s *Service) cancelTossPayment(paymentKey, reason string) error {
	log.Printf("결제 취소 API 시작")
	secretKey := os.Getenv("TOSS_SECRET_KEY")
	authHeader := "Basic " + base64.StdEncoding.EncodeToString([]byte(secretKey+":"))

	cancelData := map[string]string{"cancelReason": reason}
	jsonBody, _ := json.Marshal(cancelData)

	url := fmt.Sprintf("https://api.tosspayments.com/v1/payments/%s/cancel", paymentKey)
	req, _ := http.NewRequest("POST", url, bytes.NewBuffer(jsonBody))
	req.Header.Set("Authorization", authHeader)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		log.Printf("[토스 취소 거절 사유]: %s", string(body))
		return errors.New("토스 취소 API 호출 실패")
	}
	log.Printf("[SUCCESS] 결제 취소 완료: %s (사유: %s)", paymentKey, reason)
	return nil
}

func (s *Service) GetUserPaymentList(userID uint) ([]PaymentLog, error) {
	return s.paymentRepo.PaymentLogFindByUserID(userID)
}
