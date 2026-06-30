package payment

import "github.com/gin-gonic/gin"

func SetupRoutes(r *gin.Engine, h *Handler) {
	r.POST("/payments/confirm/:user_id/:location_id", h.ConfirmPayment)
	r.GET("payments/history/:id", h.GetPaymentHistory)
}
