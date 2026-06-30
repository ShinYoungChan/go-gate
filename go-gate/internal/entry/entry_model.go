package entry

import "time"

type AccessLog struct {
	ID         uint      `gorm:"primaryKey"`
	UserID     uint      `gorm:"index"`
	AccessedAt time.Time `gorm:"type:timestamp;not null"`
	Method     string    `gorm:"size:20;not null"`
	Result     string    `gorm:"size:100;not null"`
}
