package models

import "time"

// 이용권 종류 (membership_items)
type MembershipItem struct {
	ID           uint   `gorm:"primaryKey"`
	Title        string `gorm:"size:20;not null"`
	Type         string `gorm:"size:20;not null"`
	DurationDays int    `gorm:"not null"`
	Amount       int    `gorm:"not null"`
}

// 유저 이용권 (user_memberships)
type UserMembership struct {
	ID          uint      `gorm:"primaryKey"`
	UserID      uint      `gorm:"index"` // 조회 성능을 위한 인덱스 추가
	ItemID      uint      `gorm:"index"` // 조회 성능을 위한 인덱스 추가
	SttDt       time.Time `gorm:"type:timestamp;not null"`
	EndDt       time.Time `gorm:"type:timestamp;not null"`
	IsCountType bool      `gorm:"not null"`
	Count       int       `gorm:"not null"`
	IsValid     bool      `gorm:"default:true;not null"`
	Amount      int       `gorm:"not null"` // 결제 시점 금액 스냅샷
}
