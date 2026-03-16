package models

type Location struct {
	ID        uint    `gorm:"primaryKey"`
	PlaceName string  `gorm:"size:20;not null"`
	Lat       float64 `gorm:"type:decimal(10,8);not null"`
	Lon       float64 `gorm:"type:decimal(11,8);not null"`
}
