package payment

import "time"

type PaymentLog struct {
	ID               uint `gorm:"primaryKey"`
	UserID           uint `gorm:"index"`
	MembershipTypeID uint `gorm:"index"`

	PaymentKey string `gorm:"uniqueIndex"`
	OrderId    string `gorm:"uniqueIndex"`

	Amount int
	Status string
	Method string

	CreatedAt time.Time
	UpdatedAt time.Time
}
