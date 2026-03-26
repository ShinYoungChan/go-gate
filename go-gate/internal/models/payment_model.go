package models

import "time"

type PaymentLog struct {
	ID               uint `gorm:"primaryKey"`
	UserID           uint `gorm:"index"` // 누가 구입했는지
	MembershipTypeID uint `gorm:"index"` // 어떤 상품인지 (기존 membership_id)

	// --- 토스 연동 필수 필드 ---
	PaymentKey string `gorm:"uniqueIndex"` // 토스가 준 결제 고유 키 (환불 시 필수!)
	OrderId    string `gorm:"uniqueIndex"` // 우리가 생성한 주문 번호 (중복 결제 방지)

	// --- 결제 상태 및 금액 ---
	Amount int    // 실제 결제 금액
	Status string // READY, DONE, CANCELED, FAIL (상태 추적)
	Method string // CARD, TRANSFER 등 (결제 수단)

	// --- 시간 기록 ---
	CreatedAt time.Time // 구매 시점
	UpdatedAt time.Time
}
