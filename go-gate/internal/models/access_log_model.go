package models

import "time"

type AccessLog struct {
	ID         uint      `gorm:"primaryKey"`
	UserID     uint      `gorm:"index"` // 누구 기록인지 빠르게 찾기 위해 인덱스 추가
	AccessedAt time.Time `gorm:"type:timestamp;not null"`
	Method     string    `gorm:"size:20;not null"`  // QR, NFC, Bluetooth 등
	Result     string    `gorm:"size:100;not null"` // 성공, 실패(이유) 등
}
