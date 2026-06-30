package user

import "time"

type User struct {
	ID        uint      `gorm:"primaryKey"`
	Name      string    `gorm:"type:varchar(20);not null"`
	Email     string    `gorm:"type:varchar(100);unique;not null"`
	Password  string    `gorm:"size:255;not null"`
	CreatedAt time.Time `gorm:"type:timestamp"`
}
