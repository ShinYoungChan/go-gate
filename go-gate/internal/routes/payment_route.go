package routes

import (
	"go-gate/internal/handler"

	"github.com/gin-gonic/gin"
)

func SetupPaymentRoutes(r *gin.Engine, h *handler.PaymentHandler) {
	r.POST("/payments/confirm/:user_id/:location_id", h.ConfirmPayment)
	r.GET("payments/history/:id", h.GetPaymentHistory)
}
