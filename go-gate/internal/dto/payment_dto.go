package dto

type PaymentRequest struct {
	PaymentKey       string `json:"paymentKey" binding:"required"`
	OrderId          string `json:"orderId" binding:"required"`
	Amount           int    `json:"amount"`
	MembershipTypeID uint   `json:"membershipTypeId" binding:"required"`
}
