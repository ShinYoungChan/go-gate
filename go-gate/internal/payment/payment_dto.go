package payment

import "time"

type Request struct {
	PaymentKey       string `json:"paymentKey" binding:"required"`
	OrderId          string `json:"orderId" binding:"required"`
	Amount           int    `json:"amount"`
	MembershipTypeID uint   `json:"membershipTypeId" binding:"required"`
}

type Response struct {
	Message   string    `json:"message"`
	OrderId   string    `json:"orderId"`
	StartDate time.Time `json:"startDate"`
	EndDate   time.Time `json:"endDate"`
}
