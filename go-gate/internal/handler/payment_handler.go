package handler

import (
	"go-gate/internal/dto"
	"go-gate/internal/service"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type PaymentHandler struct {
	service *service.PaymentService
}

func NewPaymentHandler(service *service.PaymentService) *PaymentHandler {
	return &PaymentHandler{service: service}
}

func (h *PaymentHandler) ConfirmPayment(c *gin.Context) {
	// 1. 프론트엔드가 보낸 JSON 바인딩 (paymentKey, orderId, amount, membershipTypeId)
	var req dto.PaymentRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	// 2. 현재 로그인한 유저 ID 가져오기 (Context에서 추출)
	userIDStr := c.Param("id")
	userID, _ := strconv.Atoi(userIDStr)
	// 3. 서비스 호출: h.service.ApprovePayment(req, userID)
	// 추후 JWT 인증하면 c.GET 으로 변경 후 사용 예정
	result, err := h.service.ApprovePayment(req, uint(userID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	log.Printf("결제 성공 데이터: %+v", result)
	// 4. 결과 응답: 성공 시 200 OK, 실패 시 에러 코드
	c.JSON(http.StatusOK, result)
}
