package models

import "time"

type User struct {
	ID        uint      `gorm:"primaryKey"`
	Name      string    `gorm:"type:varchar(20);not null"`
	Email     string    `gorm:"type:varchar(100);unique;not null"`
	Password  string    `gorm:"size:255;not null"` // size:255 -> varchar(255) 랑 동일! size 쓰는 이유는 DB 범용성 때문
	CreatedAt time.Time `gorm:"type:timestamp"`
}
