package dto

import "time"

type PaymentRequest struct {
	PaymentKey       string `json:"paymentKey" binding:"required"`
	OrderId          string `json:"orderId" binding:"required"`
	Amount           int    `json:"amount"`
	MembershipTypeID uint   `json:"membershipTypeId" binding:"required"`
}

type PaymentResponse struct {
	Message   string    `json:"message"`
	OrderId   string    `json:"orderId"`
	StartDate time.Time `json:"startDate"`
	EndDate   time.Time `json:"endDate"`
}
