package membership

import "time"

// MembershipItem 이용권 종류 (membership_items)
type MembershipItem struct {
	ID           uint   `gorm:"primaryKey"`
	LocationID   uint   `gorm:"index"`
	Title        string `gorm:"size:20;not null"`
	Type         string `gorm:"size:20;not null"`
	DurationDays int    `gorm:"not null"`
	Amount       int    `gorm:"not null"`
}

// UserMembership 유저 이용권 (user_memberships)
type UserMembership struct {
	ID          uint      `gorm:"primaryKey"`
	UserID      uint      `gorm:"index"`
	LocationID  uint      `gorm:"index"`
	ItemID      uint      `gorm:"index"`
	SttDt       time.Time `gorm:"type:timestamp;not null"`
	EndDt       time.Time `gorm:"type:timestamp;not null"`
	IsCountType bool      `gorm:"not null"`
	Count       int       `gorm:"not null"`
	IsValid     bool      `gorm:"default:true;not null"`
	Amount      int       `gorm:"not null"`
}
