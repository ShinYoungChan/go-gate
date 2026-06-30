package payment

import (
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) ConfirmPayment(c *gin.Context) {
	var req Request
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID, _ := strconv.Atoi(c.Param("id"))
	// locationID 추후 JWT 인증 도입 후 context에서 추출 예정
	result, err := h.service.ApprovePayment(req, uint(userID), 1)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	log.Printf("결제 성공 데이터: %+v", result)
	c.JSON(http.StatusOK, result)
}

func (h *Handler) GetPaymentHistory(c *gin.Context) {
	userID, _ := strconv.Atoi(c.Param("id"))

	userPaymentLogs, err := h.service.GetUserPaymentList(uint(userID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "내역 조회 중 오류 발생"})
		return
	}

	if len(userPaymentLogs) == 0 {
		c.JSON(http.StatusOK, gin.H{
			"message": "결제 내역이 없습니다.",
			"data":    []interface{}{},
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "조회 성공!", "data": userPaymentLogs})
}
